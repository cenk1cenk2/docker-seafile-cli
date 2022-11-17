package pipe

import (
	"fmt"
	"path"

	"github.com/mitchellh/go-ps"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func HealthCheck(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health", "parent").
		SetJobWrapper(func(job Job) Job {
			return tl.JobSequence(
				job,
				tl.JobParallel(
					HealthCheckSeafDaemon(tl).Job(),
					HealthCheckStatus(tl).Job(),
				),
			)
		}).
		Set(func(t *Task[Pipe]) error {
			processes, err := ps.Processes()

			if err != nil {
				return err
			}

			for _, process := range processes {
				switch process.Executable() {
				case "seaf-cli":
					t.Pipe.Ctx.Health.SeafDaemonPID = append(t.Pipe.Ctx.Health.SeafDaemonPID, process.Pid())
					t.Log.Debugf("Seafile Daemon PID set: %w", t.Pipe.Ctx.Health.SeafDaemonPID)
				}
			}

			return nil
		})
}

func HealthCheckSeafDaemon(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health", "seaf-daemon").
		SetJobWrapper(func(job Job) Job {
			return tl.JobBackground(tl.JobLoopWithWaitAfter(job, tl.Pipe.Health.CheckInterval))
		}).
		Set(func(t *Task[Pipe]) error {
			for _, pid := range t.Pipe.Ctx.Health.SeafDaemonPID {
				process, err := ps.FindProcess(pid)

				if err != nil {
					t.Log.Debugln(err)
				}

				if process == nil {
					t.SendFatal(fmt.Errorf("Seafile Daemon process is not alive."))

					return nil
				}
			}

			t.Log.Debugf(
				"Next Seafile Daemon process health check in: %s",
				t.Pipe.Health.CheckInterval.String(),
			)

			return nil
		})
}

func HealthCheckStatus(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health", "status").
		SetJobWrapper(func(job Job) Job {
			return tl.JobBackground(tl.JobLoopWithWaitAfter(job, tl.Pipe.StatusInterval))
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"status",
				"-c",
				path.Join(t.Pipe.Seafile.DataLocation, "ccnet"),
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				return err
			}

			t.Log.Debugf(
				"Next status check in: %s",
				t.Pipe.Health.StatusInterval.String(),
			)

			return nil
		})
}
