package pipe

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Health struct {
		StatusInterval time.Duration
	}

	Server struct {
		Url string `validate:"url"`
	}

	Credentials struct {
		Username string
		Password string `validate:"required_without=Token"`
		Token    string `validate:"required_without=Password"`
	}

	Seafile struct {
		MountLocation string `validate:"dir"`
		DataLocation  string `validate:"dir"`
	}

	Pipe struct {
		Health
		Server
		Credentials
		Seafile
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).Set(
		func(tl *TaskList) Job {
			return JobSequence(
				Tasks(tl).Job(),
				Services(tl).Job(),
				HealthCheck(tl).Job(),
				JobWaitForTerminator(p),
			)
		})
}
