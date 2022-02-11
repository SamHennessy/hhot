package domain

import (
	"context"
	"os/exec"
	"syscall"

	"github.com/SamHennessy/hhot/hhotui/core"
	"github.com/SamHennessy/hhot/hhotui/domain/topics"
)

func NewApp(sl ServiceLocator) *App {
	return &App{
		sl: sl,
	}
}

type App struct {
	sl      ServiceLocator
	Build   Build
	Runtime Runtime
	cancel  context.CancelFunc
}

func (a *App) Start(ctx context.Context) {
	ctx, a.cancel = context.WithCancel(ctx)
	a.Runtime = NewRuntime(a.newCmd(ctx))

	// Hack
	a.Runtime.Status = RuntimeStatusRunning

	a.sl.AppPubSub().Publish(topics.AppStart, a)

	// Blocking call
	a.Runtime.runCmd()

	a.sl.AppPubSub().Publish(topics.AppStop, a)
}

func (a *App) Stop() error {
	if a.Runtime.Status == RuntimeStatusEmpty {
		return nil
	}

	if a.Runtime.Status == RuntimeStatusStopped {
		// <-a.Runtime.done
		return nil
	}

	a.cancel()

	var err error

	if a.Runtime.cmd != nil && a.Runtime.cmd.Process != nil {
		err = syscall.Kill(-a.Runtime.cmd.Process.Pid, syscall.SIGTERM)
	}

	return err
}

func (a *App) newCmd(ctx context.Context) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "./bin/hhot_tmp", "-base-path", "/app/")

	errOutReader := core.NewLogReader(a.sl.AppPubSub(), topics.OutputAppErr)
	stdOutReader := core.NewLogReader(a.sl.AppPubSub(), topics.OutputAppStd)

	cmd.Stderr = errOutReader
	cmd.Stdout = stdOutReader

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd
}
