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

		path: l.NewLockBox("Loading..."),
	}

	c.Add(
		l.T("input",
			l.Class("w-full px-2 py-1"),
			l.Attrs{"autocomplete": "off"},
			l.AttrsLockBox{"value": c.path},
		),
	)

	return c
}

type browserPath struct {
	*l.Component

	path *l.LockBox[string]
}

func (b *browserPath) PubSubMount(ctx context.Context, pubSub *hlivekit.PubSub) {
	pubSub.Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		data, ok := message.Value.(iframeData)

		if !ok {
			return
		}

		b.path.Set(data.path)
	}), topics.IframeUpdate)
}

func newBrowserTitle(elements ...any) *browserTitle {
	c := &browserTitle{
		Component: l.C("div", l.Class("px-2 truncate text-gray-300 font-thin"), elements),

		title: l.Box("Loading..."),
	}

	c.Add(c.title)

	return c
}

type browserTitle struct {
	*l.Component

	title *l.NodeBox[string]
}

func (b *browserTitle) PubSubMount(ctx context.Context, pubSub *hlivekit.PubSub) {
	pubSub.Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		data, ok := message.Value.(iframeData)

		if !ok {
			return
		}

		b.title.Set(data.title)
	}), topics.IframeUpdate)
}
