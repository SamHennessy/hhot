package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/SamHennessy/hhot/hhotui/domain/topics"
	l "github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
)

func logList(sl ServiceLocator) l.Tagger {
	c := hlivekit.List("div", l.Class(""))

	var subFn hlivekit.QueueSubscriber

	topicList := []string{topics.LogUI, topics.OutputAppErr, topics.OutputAppStd, topics.LogsClear, topics.LogsAddMarker}

	c.MountFunc = func(ctx context.Context) {
		subFn = hlivekit.NewSub(func(message hlivekit.QueueMessage) {
			switch message.Topic {
			case topics.LogsClear:
				c.RemoveAllItems()
			case topics.LogsAddMarker:
				line := l.CM("div", l.T("div", l.Class("py-1 my-2 bg-gradient-to-r from-gray-600 to-gray-300")))
				line.Add(
					scrollToView(),
					hlivekit.OnDiffApplyOnce(func(ctx context.Context, e l.Event) {
						line.Add(scrollToViewRemove())
					}),
				)

				c.Add(line)

			case topics.LogUI, topics.OutputAppErr, topics.OutputAppStd:
				log, ok := message.Value.([]byte)
				if !ok {
					return
				}

				var logJSON map[string]interface{}

				if err := json.Unmarshal(log, &logJSON); err == nil {
					var keys []string

					var newLog string
					for key := range logJSON {
						keys = append(keys, key)
					}

					sort.Strings(keys)

					for i := 0; i < len(keys); i++ {
						key := keys[i]
						newLog += key + "='" + fmt.Sprint(logJSON[key]) + "' "
					}

					// fmt.Println("Log Reader Error: ", err)

					// return
					log = []byte(newLog)
				}

				prefix := "[ hhot ] "
				switch message.Topic {
				case topics.OutputAppErr:
					prefix = "[stderr] "
				case topics.OutputAppStd:
					prefix = "[stdout] "
				}

				logs := strings.Split(string(log), "\n")

				// Typically, the last entry will be unwanted
				if len(logs) > 1 {
					if logs[len(logs)-1] == "" {
						logs = logs[:len(logs)-1]
					}
				}

				for i := 0; i < len(logs); i++ {
					line := l.CM("div", l.T("pre", l.Style{"margin": "0"}, prefix+logs[i]))

					// Last item
					if i == len(logs)-1 {
						line.Add(
							scrollToView(),
							hlivekit.OnDiffApplyOnce(func(ctx context.Context, e l.Event) {
								line.Add(scrollToViewRemove())
							}),
						)
					}

					c.Add(line)
				}
			}

			// This leaves some components sometimes
			// l.RenderComponent(ctx, c)
			l.Render(ctx)
		})

		sl.AppPubSub().Subscribe(subFn, topicList...)
	}

	c.UnmountFunc = func(ctx context.Context) {
		sl.AppPubSub().Unsubscribe(subFn, topicList...)
	}

	return c
}
