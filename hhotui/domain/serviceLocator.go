package domain

import (
	"github.com/SamHennessy/hlive/hlivekit"
	"github.com/rs/zerolog"
)

type ServiceLocator interface {
	AppPubSub() *hlivekit.PubSub
	Logger() *zerolog.Logger
}
