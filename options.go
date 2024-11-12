package chromedpool

import (
	"time"

	"github.com/chromedp/chromedp"
)

type Option func(*TabPool)

func WithMaxTabs(maxTabs int) Option {
	return func(tp *TabPool) {
		tp.maxTabs = maxTabs
	}
}

func WithHeadless(headless bool) Option {
	return func(tp *TabPool) {
		tp.headless = headless
	}
}

func WithUserAgent(userAgent string) Option {
	return func(tp *TabPool) {
		tp.userAgent = userAgent
	}
}

func WithProxy(proxyURL string) Option {
	return func(tp *TabPool) {
		tp.proxyURL = proxyURL
	}
}

func WithWaitTimeout(waitTimeout time.Duration) Option {
	return func(tp *TabPool) {
		tp.waitTimeout = waitTimeout
	}
}

func WithChromeFlag(name string, value interface{}) Option {
	return func(tp *TabPool) {
		tp.chromeFlags = append(tp.chromeFlags, chromedp.Flag(name, value))
	}
}
