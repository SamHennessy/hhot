package domain

import (
	"os/exec"
	"time"

	"github.com/SamHennessy/hhot/hhotui/domain/topics"
)

type BuildStatus int

const (
	BuildStatusEmpty BuildStatus = iota
	BuildStatusSuccess
	BuildStatusFailed
	BuildStatusInProgress
)

type Build struct {
	Status     BuildStatus
	Output     string
	BuildStart time.Time
	BuildEnd   time.Time
}

func (b *Build) BuildDuration() time.Duration {
	if b.BuildStart.IsZero() {
		return 0
	}

	end := b.BuildEnd
	if end.IsZero() {
		end = time.Now()
	}

	return end.Sub(b.BuildStart)
}

func StartBuild(sl ServiceLocator) error {
	b := Build{}
	b.BuildStart = time.Now()
	cmd := exec.Command("go", "build", "-o", "./bin/hhot_tmp", "./cmd/server/main.go")

	sl.Logger().Trace().Str("cmd", cmd.String()).Msg("build command")

	b.Status = BuildStatusInProgress

	sl.AppPubSub().Publish(topics.Build, b)

	out, err := cmd.CombinedOutput()

	b.BuildEnd = time.Now()
	b.Output = string(out)
	b.Status = BuildStatusSuccess
	if err != nil {
		b.Status = BuildStatusFailed
	}

	sl.AppPubSub().Publish(topics.Build, b)

	return err
}
