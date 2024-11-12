package chromedpool

import "time"

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

func WithChromeFlags(flags ...string) Option {
	return func(tp *TabPool) {
		tp.chromeFlags = append(tp.chromeFlags, flags...)
	}
}
