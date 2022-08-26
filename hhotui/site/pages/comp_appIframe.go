package pages

import (
	"context"
	"sync"

	"github.com/SamHennessy/hhot/hhotui/domain"
	"github.com/SamHennessy/hhot/hhotui/domain/topics"
	l "github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
)

func newAppIframe(sl ServiceLocator) *appIframe {
	frame := l.C("div", l.Class("text-4xl animate-pulse"), "Loading...")

	c := &appIframe{
		frame:              frame,
		app:                sl.App(),
		ComponentMountable: l.CM("div", l.Class("w-full h-full"), frame),
	}

	var subFn hlivekit.QueueSubscriber

	c.MountFunc = func(ctx context.Context) {
		subFn = hlivekit.NewSub(func(message hlivekit.QueueMessage) {
			if message.Topic == topics.RefreshDo {
				IframeReload(frame)
			} else {
				c.render()
			}

			l.Render(ctx)
		})

		// Listen for updates
		sl.AppPubSub().Subscribe(subFn, topics.AppStart, topics.AppStop, topics.RefreshDo)

		// Init
		c.render()
	}

	c.UnmountFunc = func(ctx context.Context) {
		sl.AppPubSub().Unsubscribe(subFn, topics.AppStart, topics.AppStop, topics.RefreshDo)
	}

	return c
}

type appIframe struct {
	app    *domain.App
	once   sync.Once
	frame  *l.Component
	pubSub *hlivekit.PubSub

	*l.ComponentMountable
}

type iframeData struct {
	title string
	path  string
}

func (c *appIframe) render() {
	if c.app.Runtime.Status != domain.RuntimeStatusRunning {
		c.Add(l.ClassBool{"hidden": true})
	} else {
		c.Add(l.ClassBool{"hidden": false})

		var first bool

		c.once.Do(func() {
			*c.frame = *l.C("iframe",
				l.Class("w-full h-full"),
				l.Attrs{"src": "http://localhost:3333/app"},
				InstallIframeWatcher(func(ctx context.Context, e l.Event) {
					c.pubSub.Publish(topics.IframeUpdate, iframeData{
						title: e.Extra["title"],
						path:  e.Extra["path"],
					})
				}),
			)

			first = true
		})

		if !first {
			IframeReload(c.frame)
		}
	}
}

func (c *appIframe) PubSubMount(_ context.Context, pubSub *hlivekit.PubSub) {
	c.pubSub = pubSub
}
