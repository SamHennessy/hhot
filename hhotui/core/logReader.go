package core

import (
	"github.com/SamHennessy/hlive/hlivekit"
)

func NewLogReader(pubSub *hlivekit.PubSub, topic string) *LogReader {
	l := &LogReader{
		pubSub: pubSub,
		topic:  topic,
	}

	return l
}

type LogReader struct {
	pubSub *hlivekit.PubSub
	topic  string
}

func (l *LogReader) Write(p []byte) (int, error) {
	l.pubSub.Publish(l.topic, p)

	return len(p), nil
}
