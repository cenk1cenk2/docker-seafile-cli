package pipe

import (
	"fmt"

	"github.com/mitchellh/go-ps"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func HealthCheck(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health", "parent").
		SetJobWrapper(func(job Job) Job {
			return TL.JobSequence(
				job,
				TL.JobParallel(
					HealthCheckSeafDaemon(tl).Job(),
					HealthCheckCcnet(tl).Job(),
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
				case "seaf-daemon":
					t.Pipe.Ctx.Health.SeafDaemonPID = process.Pid()
					t.Log.Debugf("Seafile Daemon PID set: %d", t.Pipe.Ctx.Health.SeafDaemonPID)
				case "ccnet":
					t.Pipe.Ctx.Health.CcnetPID = process.Pid()
					t.Log.Debugf("CCNet PID set: %d", t.Pipe.Ctx.Health.CcnetPID)
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
			process, err := ps.FindProcess(t.Pipe.Ctx.Health.SeafDaemonPID)

			if err != nil {
				t.Log.Debugln(err)
			}

			if process == nil {
				t.SendFatal(fmt.Errorf("Seafile Daemon process is not alive."))

				return nil
			}

			t.Log.Debugf(
				"Next Seafile Daemon process health check in: %s",
				t.Pipe.Health.CheckInterval.String(),
			)

			return nil
		})
}

func HealthCheckCcnet(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("health", "ccnet").
		SetJobWrapper(func(job Job) Job {
			return tl.JobBackground(tl.JobLoopWithWaitAfter(job, tl.Pipe.Health.CheckInterval))
		}).
		Set(func(t *Task[Pipe]) error {
			process, err := ps.FindProcess(t.Pipe.Ctx.Health.CcnetPID)

			if err != nil {
				t.Log.Debugln(err)
			}

			if process == nil {
				t.SendFatal(fmt.Errorf("CCNet process is not alive."))

				return nil
			}

			t.Log.Debugf(
				"Next CCNet process health check in: %s",
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
			).
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