package site

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/SamHennessy/hhot/hhotui/domain"
	"github.com/SamHennessy/hhot/hhotui/domain/topics"
	"github.com/SamHennessy/hlive/hlivekit"
	"github.com/fsnotify/fsnotify"
)

func Serve(sl *ServiceLocatorSite) {
	if sl == nil {
		sl = NewServiceLocatorSite()
	}

	Router(sl)

	sl.Logger().Info().Str("address", "http://localhost:3333").Msg("start ui")
	log.Println("HHot UI: http://localhost:3333")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Println("HHot UI: ERROR: http listen and serve: ", err)
	}
}

func NewServer(sl *ServiceLocatorSite) *Server {
	return &Server{sl: sl}
}

type Server struct {
	sl *ServiceLocatorSite
}

func (s *Server) MainLoop(outerCTX context.Context, watcher *fsnotify.Watcher) {
	type lastTrigger struct {
		when time.Time
		name string
	}

	var (
		ctx    context.Context
		cancel context.CancelFunc

		cmdExit           = make(chan bool)
		triggerBuild      = make(chan bool)
		triggerBuildJS    = make(chan bool)
		triggerBuildCSS   = make(chan bool)
		triggerAppStop    = make(chan bool)
		triggerAppStart   = make(chan bool)
		triggerAppRestart = make(chan bool)

		lastTriggerBuild lastTrigger
	)

	go domain.AssetBuild()

	s.sl.Logger().Info().Msg("build app")
	err := domain.StartBuild(s.sl)
	if err != nil {
		s.sl.Logger().Err(err).Msg("build")
	} else {
		// Start App
		ctx, cancel = context.WithCancel(context.Background())

		go func() {
			s.sl.Logger().Info().Msg("start app")
			s.sl.App().Start(ctx)
			cmdExit <- true
		}()
	}

	s.sl.AppPubSub().Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		triggerBuild <- true
	}), topics.TriggerBuild)

	s.sl.AppPubSub().Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		triggerBuildJS <- true
	}), topics.TriggerBuildJS)

	s.sl.AppPubSub().Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		triggerBuildCSS <- true
	}), topics.TriggerBuildCSS)

	s.sl.AppPubSub().Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		triggerAppStop <- true
	}), topics.AppStopDo)

	s.sl.AppPubSub().Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		triggerAppStart <- true
	}), topics.AppStartDo)

	s.sl.AppPubSub().Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		triggerAppRestart <- true
	}), topics.AppRestartDo)

	loop := true
	for loop {
		select {
		// TODO: create a trigger flag, to prevent queuing of builds
		case <-triggerBuild:
			assetBuildDone := make(chan bool)
			go func() {
				domain.AssetBuild()
				assetBuildDone <- true
			}()

			s.sl.Logger().Info().Msg("stop app")

			if cancel != nil {
				cancel()
			}

			if err := s.sl.App().Stop(); err != nil {
				s.sl.Logger().Err(err).Msg("stop app")
			}

			s.sl.Logger().Info().Msg("build app")

			if err := domain.StartBuild(s.sl); err != nil {
				s.sl.Logger().Err(err).Msg("build app")

				<-assetBuildDone
			} else {
				<-assetBuildDone

				ctx, cancel = context.WithCancel(context.Background())

				go func() {
					s.sl.Logger().Info().Msg("start app")
					s.sl.App().Start(ctx)
					cmdExit <- true
				}()
			}
		case <-triggerBuildJS:
			go domain.AssetBuild()
		case <-triggerBuildCSS:
			go domain.AssetBuild()
		case event, ok := <-watcher.Events:
			if !ok {
				loop = false
				s.sl.Logger().Error().Msg("watcher events channel closed")
			}

			s.sl.Logger().Trace().Str("name", event.Name).Str("op", event.Op.String()).Msg("watcher event")
			if event.Op&fsnotify.Write == fsnotify.Write {

				if strings.HasSuffix(event.Name, ".go") {
					now := time.Now()

					if now.Sub(lastTriggerBuild.when) > 5*time.Second {
						lastTriggerBuild.when = now
						lastTriggerBuild.name = event.Name

						go s.sl.AppPubSub().Publish(topics.TriggerBuild, nil)
					}
				} else if strings.HasSuffix(event.Name, ".json") || strings.HasSuffix(event.Name, ".js") {
					go s.sl.AppPubSub().Publish(topics.TriggerBuildJS, nil)
				} else if strings.HasSuffix(event.Name, ".css") {
					go s.sl.AppPubSub().Publish(topics.TriggerBuildCSS, nil)
				}
			}

			// Watch new folders
			if event.Op&fsnotify.Create == fsnotify.Create {
				file, err := os.Stat(event.Name)
				if err != nil {
					s.sl.Logger().Err(err).Msg("os.Stat")
				} else if strings.HasSuffix(file.Name(), "~") {
					// Nop
				} else if file.IsDir() && (file.Name() != "dist" || file.Name() != "node_modules") {
					if err := watcher.Add(event.Name); err != nil {
						s.sl.Logger().Err(err).Msg("ERROR: watch new dir")
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				loop = false
				s.sl.Logger().Error().Msg("watcher error channel closed")
			}
			s.sl.Logger().Err(err).Msg("watcher")

		case <-triggerAppStop:
			if err := s.sl.App().Stop(); err != nil {
				s.sl.Logger().Err(err).Msg("stop app")
			}
		case <-triggerAppStart:
			if s.sl.App().Build.Status == domain.BuildStatusSuccess && s.sl.App().Runtime.Status == domain.RuntimeStatusStopped {
				go func() {
					s.sl.Logger().Info().Msg("start app")
					s.sl.App().Start(ctx)
					cmdExit <- true
				}()
			}
		case <-triggerAppRestart:
			if err := s.sl.App().Stop(); err != nil {
				s.sl.Logger().Err(err).Msg("stop app")
			}

			if s.sl.App().Build.Status == domain.BuildStatusSuccess {

				for i := 0; s.sl.App().Runtime.Status == domain.RuntimeStatusRunning && i < 5; i++ {
					time.Sleep(100 * time.Millisecond)
				}

				if s.sl.App().Runtime.Status == domain.RuntimeStatusStopped {
					go func() {
						s.sl.Logger().Info().Msg("start app")
						s.sl.App().Start(ctx)
						cmdExit <- true
					}()
				}
			}
		case <-cmdExit:
			s.sl.Logger().Info().Msg("app stopped")
		case <-outerCTX.Done():
			loop = false
		}
	}

	s.sl.Logger().Trace().Msg("main loop end")

	if cancel != nil {
		cancel()
	}
}
