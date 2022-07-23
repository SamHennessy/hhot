package pages

import (
	_ "embed"

	l "github.com/SamHennessy/hlive"
)

const (
	onLoadAttrName = "hhui-onload-url"
)

func onloadurl() l.AttributePluginer {
	return &OnloadAttribute{
		Attribute: l.NewAttribute(onLoadAttrName, ""),
	}
}

type OnloadAttribute struct {
	*l.Attribute
}

//go:embed attr_onloadurl.js
var attrOnloadURL []byte

func (a *OnloadAttribute) Initialize(page *l.Page) {
	page.DOM.Head.Add(l.T("script", l.HTML(attrOnloadURL)))
}

func (a *OnloadAttribute) InitializeSSR(page *l.Page) {
}
