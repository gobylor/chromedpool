package chromedpool

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

type TabPool struct {
	pool     *sync.Pool
	browser  context.Context
	cancel   context.CancelFunc
	tabCount int
	mu       sync.Mutex

	maxTabs     int
	headless    bool
	userAgent   string
	proxyURL    string
	chromeFlags []string
}

func NewTabPool(options ...Option) (*TabPool, error) {
	tp := &TabPool{
		maxTabs:  10,
		headless: true,
	}

	for _, option := range options {
		option(tp)
	}

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
		opts = append(opts, chromedp.Flag(flag, true))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	browserCtx, _ := chromedp.NewContext(allocCtx)
	err := chromedp.Run(browserCtx,
		chromedp.Navigate("about:blank"),
	)
	if err != nil {
		return nil, err
	}
	tp.browser = browserCtx
	tp.cancel = cancel

	tp.pool = &sync.Pool{
		New: func() interface{} {
			tp.mu.Lock()
			defer tp.mu.Unlock()
			if tp.tabCount >= tp.maxTabs {
				return nil
			}
			ctx, cancel := chromedp.NewContext(tp.browser)
			tp.tabCount++
			return Tab{
				ctx:    ctx,
				cancel: cancel,
			}
		},
	}
	return tp, nil
}

func (tp *TabPool) GetTab() Tab {
	tab := tp.pool.Get()
	if tab == nil {
		time.Sleep(100 * time.Millisecond)
		return tp.GetTab()
	}
	return tab.(Tab)
}

func (tp *TabPool) PutTab(tab Tab) {
	err := chromedp.Run(tab.ctx, chromedp.Navigate("about:blank"))
	if err != nil {
		log.Printf("Error clearing tab: %v", ErrPoolClosed)
	}
	tp.pool.Put(tab)
}

func (tp *TabPool) Close() {
	tp.cancel()
}
