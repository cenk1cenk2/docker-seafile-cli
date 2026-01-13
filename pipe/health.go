package pipe

import (
	"fmt"
	"path"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/mitchellh/go-ps"
)

func HealthCheck(tl *TaskList) *Task {
	return tl.CreateTask("health", "parent").
		SetJobWrapper(func(job Job, _ *Task) Job {
			return JobParallel(
				// HealthCheckSeafDaemon(tl).Job(),
				HealthCheckStatus(tl).Job(),
			)
		})
}

func HealthCheckSeafDaemon(tl *TaskList) *Task {
	return tl.CreateTask("health", "seaf-daemon").
		SetJobWrapper(func(job Job, _ *Task) Job {
			return JobLoopWithWaitAfter(job, P.Health.CheckInterval)
		}).
		Set(func(t *Task) error {
			if len(C.pids) == 0 {
				return (fmt.Errorf("Seafile Daemon process PID not found."))
			}

			for _, pid := range C.pids {
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
				P.Health.CheckInterval.String(),
			)

			return nil
		})
}

func HealthCheckStatus(tl *TaskList) *Task {
	return tl.CreateTask("health", "status").
		SetJobWrapper(func(job Job, _ *Task) Job {
			return JobLoopWithWaitAfter(job, P.Health.StatusInterval)
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"status",
				"-c",
				path.Join(P.Seafile.DataLocation, "ccnet"),
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				return err
			}

			t.Log.Debugf(
				"Next status check in: %s",
				P.Health.StatusInterval.String(),
			)

			return nil
		})
}
