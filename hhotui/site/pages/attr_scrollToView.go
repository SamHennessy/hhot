package pages

import (
	_ "embed"

	l "github.com/SamHennessy/hlive"
)

const (
	scrollToViewAttrName = "hhui-scrollToView"
)

func scrollToView() l.AttributePluginer {
	return &scrollToViewAttribute{
		Attribute: l.NewAttribute(scrollToViewAttrName, ""),
	}
}

func scrollToViewRemove() l.AttrsOff {
	return l.AttrsOff{scrollToViewAttrName}
}

type scrollToViewAttribute struct {
	*l.Attribute
}

//go:embed attr_scrollToView.js
var scrollToViewJS []byte

func (a *scrollToViewAttribute) Initialize(page *l.Page) {
	page.DOM().Head().Add(l.T("script", l.HTML(scrollToViewJS)))
}

func (a *scrollToViewAttribute) InitializeSSR(page *l.Page) {
}
