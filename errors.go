package chromedpool

import "errors"

var (
	ErrPoolClosed = errors.New("tab pool is closed")
	ErrTimeout    = errors.New("timeout waiting for available tab")
)
