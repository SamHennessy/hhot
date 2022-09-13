package hhot

import (
	"context"
	_ "embed"
	"strings"

	l "github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
)

func NewPageHistoryManager(config Config) *PageHistoryManager {
	phm := &PageHistoryManager{
		config: config,
	}
	phm.eb = l.On("popstate", phm.OnPopState)
	phm.eb.ID = "phm"

	return phm
}

func (phm *PageHistoryManager) OnPopState(_ context.Context, e l.Event) {
	phm.pubSub.Publish(TopicRedirectInternalHistory, e.Extra["path"])
}

func (phm *PageHistoryManager) InstallPageHistory(pubSub *hlivekit.PubSub) *PageHistory {
	phm.pubSub = pubSub
	a := &PageHistory{
		Attribute: l.NewAttribute(pageHistoryAttrNameOnPopState, ""),
		pubSub:    pubSub,
		eb:        phm.eb,
		config:    phm.config,
	}

	return a
}

type PageHistoryManager struct {
	config Config
	pubSub *hlivekit.PubSub
	attr   *PageHistory
	eb     *l.EventBinding
}

const (
	pageHistoryAttrNameOnPopState        = "hh_onpopstate"
	pageHistoryAttrNamePush              = "hh_history_push"
	pageHistoryEventBindingIDTemplateVar = "__bindingID__"
	pageHistoryEventAttrTemplateVar      = "__pushAttrName__"
	pageHistoryTemplateVarBasePath       = "__base_path__"
)

type PageHistory struct {
	*l.Attribute

	eb       *l.EventBinding
	pubSub   *hlivekit.PubSub
	config   Config
	rendered bool
}

//go:embed history.js
var historyJS []byte

func (a *PageHistory) Initialize(page *l.Page) {
	if a.rendered {
		return
	}

	jsStr := strings.ReplaceAll(string(historyJS), pageHistoryEventBindingIDTemplateVar, a.eb.ID)
	jsStr = strings.ReplaceAll(jsStr, pageHistoryEventAttrTemplateVar, pageHistoryAttrNamePush)
	jsStr = strings.ReplaceAll(jsStr, pageHistoryTemplateVarBasePath, a.config.BasePath())

	page.DOM().Head().Add(l.T("script", l.HTML(jsStr)), a.eb)
}

func (a *PageHistory) InitializeSSR(page *l.Page) {
	a.rendered = true

	jsStr := strings.ReplaceAll(string(historyJS), pageHistoryEventBindingIDTemplateVar, a.eb.ID)
	jsStr = strings.ReplaceAll(jsStr, pageHistoryEventAttrTemplateVar, pageHistoryAttrNamePush)
	jsStr = strings.ReplaceAll(jsStr, pageHistoryTemplateVarBasePath, a.config.BasePath())

	page.DOM().Head().Add(l.T("script", l.HTML(jsStr)), a.eb)
}

func HistoryPush(path string, c l.Adder) {
	c.Add(l.Attrs{pageHistoryAttrNamePush: path})
	c.Add(hlivekit.OnDiffApplyOnce(func(ctx context.Context, e l.Event) {
		c.Add(l.Attrs{pageHistoryAttrNamePush: nil})
	}))
}
