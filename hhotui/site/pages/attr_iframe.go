package pages

import (
	"context"
	_ "embed"
	"strings"

	l "github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
)

const (
	iframeAttrReload   = "hhui-iframe-reload"
	iframeAttrWatcher  = "hhui-iframe-watch"
	iframeWatcherEvent = "iframe-change"
	iframeWatcherDelay = "500"
)

func IframeReload(c l.Adder) {
	c.Add(newIframeReload())
	c.Add(hlivekit.OnDiffApplyOnce(func(ctx context.Context, e l.Event) {
		c.Add(l.Attrs{iframeAttrReload: nil})
	}))
}

func newIframeReload() *iframeReload {
	return &iframeReload{
		Attribute: l.NewAttribute(iframeAttrReload, ""),
	}
}

type iframeReload struct {
	*l.Attribute
}

//go:embed attr_iframeReload.js
var iframeReloadJS []byte

func (a *iframeReload) Initialize(page *l.Page) {
	js := strings.ReplaceAll(string(iframeReloadJS), "__iframeAttrReload__", iframeAttrReload)
	page.DOM.Head.Add(l.T("script", l.HTML(js)))
}

func (a *iframeReload) InitializeSSR(page *l.Page) {

}

type iframeWatcher struct {
	*l.Attribute
}

func InstallIframeWatcher(eventHandler l.EventHandler) *l.ElementGroup {
	return l.Elements(
		&iframeWatcher{l.NewAttribute(iframeAttrWatcher, "")},
		l.On(iframeWatcherEvent, eventHandler),
	)
}

//go:embed attr_iframeWatcher.js
var iframeWatcherJS []byte

func (a *iframeWatcher) Initialize(page *l.Page) {
	js := strings.ReplaceAll(string(iframeWatcherJS), "__iframeWatcherEvent__", iframeWatcherEvent)
	js = strings.ReplaceAll(js, "__iframeWatcherDelay__", iframeWatcherDelay)
	page.DOM.Head.Add(l.T("script", l.HTML(js)))
}

func (a *iframeWatcher) InitializeSSR(page *l.Page) {
	// Nop
}
