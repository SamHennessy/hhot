package domain

import (
	"os/exec"
	"time"
)

type RuntimeStatus int

const (
	RuntimeStatusEmpty = iota
	RuntimeStatusRunning
	RuntimeStatusStopped
)

func NewRuntime(cmd *exec.Cmd) Runtime {
	return Runtime{
		cmd: cmd,
	}
}

type Runtime struct {
	cmd    *exec.Cmd
	Error  error
	Status RuntimeStatus
	start  time.Time
}

func (r *Runtime) runCmd() {
	r.Status = RuntimeStatusRunning
	r.start = time.Now()

	err := r.cmd.Run()
	if err != nil {
		r.Error = err
	}

	r.Status = RuntimeStatusStopped
}

func (r *Runtime) Time() time.Duration {
	if r.start.IsZero() {
		return 0
	}

	return time.Since(r.start)
}
