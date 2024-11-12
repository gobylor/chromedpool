package chromedpool

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chromedp/chromedp"
)

type TabPool struct {
	pool     *sync.Pool
	browser  context.Context
	cancel   context.CancelFunc
	tabCount int32

	tabChan   chan *Tab
	closeChan chan struct{}

	maxTabs     int
	headless    bool
	userAgent   string
	proxyURL    string
	waitTimeout time.Duration
	chromeFlags []chromedp.ExecAllocatorOption
}

type ChromeFlag struct {
	name  string
	value interface{}
}

func NewTabPool(options ...Option) (*TabPool, error) {
	tp := &TabPool{
		maxTabs:  10,
		headless: true,
	}

	for _, option := range options {
		option(tp)
	}

	tp.tabChan = make(chan *Tab, tp.maxTabs)
	tp.closeChan = make(chan struct{})

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", tp.headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	if tp.userAgent != "" {
		opts = append(opts, chromedp.UserAgent(tp.userAgent))
	}
	if tp.proxyURL != "" {
		opts = append(opts, chromedp.ProxyServer(tp.proxyURL))
	}
	for _, flag := range tp.chromeFlags {
		opts = append(opts, flag)
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, _ := chromedp.NewContext(allocCtx)
	err := chromedp.Run(browserCtx,
		chromedp.Navigate("about:blank"),
	)
	if err != nil {
		cancel()
		return nil, err
	}
	tp.browser = browserCtx
	tp.cancel = cancel

	tp.pool = &sync.Pool{
		New: func() interface{} {
			if atomic.LoadInt32(&tp.tabCount) >= int32(tp.maxTabs) {
				return nil
			}
			tab := NewTab(tp.browser)
			atomic.AddInt32(&tp.tabCount, 1)
			return tab
		},
	}

	go tp.manageTabAvailability()

	return tp, nil
}

func (tp *TabPool) manageTabAvailability() {
	for {
		select {
		case <-tp.closeChan:
			return
		default:
			tab := tp.pool.Get()
			if tab != nil {
				select {
				case tp.tabChan <- tab.(*Tab):
				case <-tp.closeChan:
					return
				}
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (tp *TabPool) GetTab() (*Tab, error) {
	if tp.waitTimeout == 0 {
		select {
		case tab := <-tp.tabChan:
			return tab, nil
		case <-tp.closeChan:
			return nil, ErrPoolClosed
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), tp.waitTimeout)
	defer cancel()
	select {
	case tab := <-tp.tabChan:
		return tab, nil
	case <-ctx.Done():
		return nil, ErrTimeout
	case <-tp.closeChan:
		return nil, ErrPoolClosed
	}
}

func (tp *TabPool) PutTab(tab *Tab) {
	err := chromedp.Run(tab.ctx, chromedp.Navigate("about:blank"))
	if err != nil {
		log.Printf("Error clearing tab: %v", ErrPoolClosed)
		tab.cancel()
		atomic.AddInt32(&tp.tabCount, -1)
		return
	}
	tp.pool.Put(tab)
}

func (tp *TabPool) Close() {
	close(tp.closeChan)
	tp.cancel()
	for {
		select {
		case tab := <-tp.tabChan:
			tab.cancel()
		default:
			return
		}
	}
}

func (tp *TabPool) Run(tasks ...chromedp.Action) error {
	tab, err := tp.GetTab()
	if err != nil {
		return err
	}
	defer tp.PutTab(tab)
	return tab.Run(tasks...)
}

func (tp *TabPool) Navigate(url string) error {
	tab, err := tp.GetTab()
	if err != nil {
		return err
	}
	defer tp.PutTab(tab)
	return tab.Navigate(url)
}
