package hhot

import (
	"github.com/SamHennessy/hlive"
	"github.com/rs/zerolog"
)

type ServiceLocator interface {
	Config() Config
	Logger() *zerolog.Logger
	PageSessionStore() *hlive.PageSessionStore
}

//
// func NewServiceLocatorSite() *ServiceLocatorSite {
// 	sl := &ServiceLocatorSite{
// 		ServiceLocator:   core.NewServiceLocator(),
// 		config:           &Config{},
// 		pageSessionStore: hlive.NewPageSessionStore(),
// 	}
//
// 	sl.app = domain.NewApp(sl)
//
// 	return sl
// }
//
// type ServiceLocatorSite struct {
// 	*core.ServiceLocator
//
// 	app              *domain.App
// 	config           *Config
// 	pageSessionStore *hlive.PageSessionStore
// }
//
// func (sl *ServiceLocatorSite) Logger() *zerolog.Logger {
// 	return sl.ServiceLocator.Logger
// }
//
// func (sl *ServiceLocatorSite) AppPubSub() *hlivekit.PubSub {
// 	return sl.ServiceLocator.PubSub
// }
//
// func (sl *ServiceLocatorSite) App() *domain.App {
// 	return sl.app
// }
//
// func (sl *ServiceLocatorSite) Config() hhot.Config {
// 	return sl.config
// }
//
// func (sl *ServiceLocatorSite) PageSessionStore() *hlive.PageSessionStore {
// 	return sl.pageSessionStore
// }
