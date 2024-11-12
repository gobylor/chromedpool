package chromedpool

import "errors"

var (
	ErrPoolClosed = errors.New("tab pool is closed")
)
