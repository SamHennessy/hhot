package core

import (
	"github.com/SamHennessy/hhot/hhotui/domain/topics"
	"github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
	"github.com/rs/zerolog"
)

// TODO: autogen

func NewServiceLocator() *ServiceLocator {
	pubSub := hlivekit.NewPubSub()
	sl := &ServiceLocator{
		Config:           &Config{},
		PageSessionStore: hlive.NewPageSessionStore(),
		// Use Wire to help with this?
		LogReader: NewLogReader(pubSub, topics.LogUI),
		PubSub:    pubSub,
	}

	sl.Logger = zerolog.New(sl.LogReader).With().Timestamp().Logger().Level(zerolog.DebugLevel)

	return sl
}

type ServiceLocator struct {
	Config           *Config
	PageSessionStore *hlive.PageSessionStore
	Logger           zerolog.Logger
	LogReader        *LogReader
	PubSub           *hlivekit.PubSub
}
