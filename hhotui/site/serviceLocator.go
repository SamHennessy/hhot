package site

import (
	"github.com/SamHennessy/hhot"
	"github.com/SamHennessy/hhot/hhotui/core"
	"github.com/SamHennessy/hhot/hhotui/domain"
	"github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
	"github.com/rs/zerolog"
)

func NewServiceLocatorSite() *ServiceLocatorSite {
	sl := &ServiceLocatorSite{
		ServiceLocator:   core.NewServiceLocator(),
		config:           &Config{},
		pageSessionStore: hlive.NewPageSessionStore(),
	}

	sl.app = domain.NewApp(sl)

	return sl
}

type ServiceLocatorSite struct {
	*core.ServiceLocator

	app              *domain.App
	config           *Config
	pageSessionStore *hlive.PageSessionStore
}

func (sl *ServiceLocatorSite) Logger() *zerolog.Logger {
	return &sl.ServiceLocator.Logger
}

func (sl *ServiceLocatorSite) AppPubSub() *hlivekit.PubSub {
	return sl.ServiceLocator.PubSub
}

func (sl *ServiceLocatorSite) App() *domain.App {
	return sl.app
}

func (sl *ServiceLocatorSite) Config() hhot.Config {
	return sl.config
}

func (sl *ServiceLocatorSite) PageSessionStore() *hlive.PageSessionStore {
	return sl.pageSessionStore
}
