package pipe

import (
	"time"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Health struct {
		CheckInterval  time.Duration
		StatusInterval time.Duration
	}

	Server struct {
		Url string `validate:"url"`
	}

	Credentials struct {
		Username string
		Password string
	}

	Seafile struct {
		MountLocation string `validate:"dir"`
		DataLocation  string `validate:"dir"`
	}

	Pipe struct {
		Ctx

		Health
		Server
		Credentials
		Seafile
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).Set(
		func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				Tasks(tl).Job(),
				Services(tl).Job(),
				HealthCheck(tl).Job(),
				tl.JobWaitForTerminator(),
			)
		})
}
