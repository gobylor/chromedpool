package chromedpool

import (
	"context"

	"github.com/chromedp/chromedp"
)

type Tab struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTab(ctx context.Context) *Tab {
	ctx, cancel := chromedp.NewContext(ctx)
	return &Tab{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *Tab) Run(actions ...chromedp.Action) error {
	return chromedp.Run(t.ctx, actions...)
}

func (t *Tab) Navigate(url string) error {
	return chromedp.Run(t.ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"))
}

func (t *Tab) Context() context.Context {
	return t.ctx
}
