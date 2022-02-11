package pages

import (
	"github.com/SamHennessy/hhot/hhotui/domain"
	"github.com/SamHennessy/hlive/hlivekit"
)

type ServiceLocator interface {
	AppPubSub() *hlivekit.PubSub
	App() *domain.App
}
