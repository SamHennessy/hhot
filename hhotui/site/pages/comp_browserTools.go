package pages

import (
	"context"

	"github.com/SamHennessy/hhot/hhotui/domain/topics"
	l "github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
)

func newBrowserPath(elements ...interface{}) *browserPath {
	c := &browserPath{
		Component: l.C("div", elements),

		path: "Loading...",
	}

	c.Add(
		l.T("input",
			l.Class("w-full px-2 py-1"),
			l.Attrs{"value": &c.path, "autocomplete": "off"},
		),
	)

	return c
}

type browserPath struct {
	*l.Component

	path string
}

func (b *browserPath) PubSubMount(ctx context.Context, pubSub *hlivekit.PubSub) {
	pubSub.Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		data, ok := message.Value.(iframeData)

		if !ok {
			return
		}

		b.path = data.path
	}), topics.IframeUpdate)
}

func newBrowserTitle(elements ...any) *browserTitle {
	c := &browserTitle{
		Component: l.C("div", l.Class("px-2 truncate text-gray-300 font-thin"), elements),

		title: "Loading...",
	}

	c.Add(&c.title)

	return c
}

type browserTitle struct {
	*l.Component

	title string
}

func (b *browserTitle) PubSubMount(ctx context.Context, pubSub *hlivekit.PubSub) {
	pubSub.Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		data, ok := message.Value.(iframeData)

		if !ok {
			return
		}

		b.title = data.title
	}), topics.IframeUpdate)
}
