package main

import (
	"context"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/SamHennessy/hhot/hhotui/domain"
	"github.com/SamHennessy/hhot/hhotui/domain/topics"
	"github.com/SamHennessy/hhot/hhotui/site"
	"github.com/SamHennessy/hlive/hlivekit"
	"github.com/fsnotify/fsnotify"
)

func main() {
	StartHot()
}

func StartHot() {
	var (
		sl     = site.NewServiceLocatorSite()
		server = site.NewServer(sl)
		done   = make(chan bool)
	)

	// TODO: move
	sl.AppPubSub().Subscribe(hlivekit.NewSub(func(message hlivekit.QueueMessage) {
		build, ok := message.Value.(domain.Build)
		if !ok {
			return
		}

		sl.App().Build = build
	}), topics.Build)

	// Start HHot UI
	go StartHHotUI(sl)

	// Create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		sl.Logger().Fatal().Err(err).Msg("new watcher")
	}

	// Start main loop
	ctx, cancel := context.WithCancel(context.Background())
	go server.MainLoop(ctx, watcher)

	// Walk file system
	err = filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && !strings.Contains(path, "hhot") && !strings.Contains(path, "node_modules") && !strings.Contains(path, "dist") {
			if err := watcher.Add(path); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		sl.Logger().Err(err).Msg("watcher file path walk")
	}

	// Wait for term signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs

		sl.Logger().Debug().Str("signal", sig.String()).Msg("signal received")

		done <- true
	}()

	<-done

	cancel()

	sl.Logger().Debug().Msg("stop file watcher")

	if err := watcher.Close(); err != nil {
		sl.Logger().Err(err).Msg("close watcher")
	}

	if sl.App().Runtime.Status == domain.RuntimeStatusRunning {
		sl.Logger().Info().Msg("stop app")

		if err := sl.App().Stop(); err != nil {
			sl.Logger().Err(err).Msg("stop app")
		}
	}
}

func StartHHotUI(sl *site.ServiceLocatorSite) {
	site.Serve(sl)
}
