package pipe

import (
	"encoding/json"
	"fmt"
	"path"
	"slices"
	"strings"
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
)

func HealthCheck(tl *TaskList) *Task {
	return tl.CreateTask("health", "parent").
		SetJobWrapper(func(job Job, _ *Task) Job {
			return JobBackground(
				// BUG: something wrong with context finisihing early on the plumber side
				GuardAlways(
					JobDelay(
						JobLoopWithWaitAfter(
							JobParallel(
								HealthCheckStatus(tl).Job(),
								HealthCheckRepositories(tl).Job(),
							),
							P.Health.StatusInterval,
						),
						15*time.Second,
					),
				),
			)
		})
}

func HealthCheckStatus(tl *TaskList) *Task {
	return tl.CreateTask("health", "status").
		Set(func(t *Task) error {
			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"status",
				"-c",
				path.Join(P.Seafile.DataLocation, "ccnet"),
			).
				SetLogLevel(LOG_LEVEL_INFO, LOG_LEVEL_WARN, LOG_LEVEL_DEBUG).
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

func HealthCheckRepositories(tl *TaskList) *Task {
	return tl.CreateTask("health", "repositories").
		Set(func(t *Task) error {
			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"list",
				"-c",
				path.Join(P.Seafile.DataLocation, "ccnet"),
				"--json",
			).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command) error {
					var libraries []SeafCliList
					if err := json.Unmarshal([]byte(strings.Join(c.GetCombinedStream(), "\n")), &libraries); err != nil {
						return fmt.Errorf("failed to parse seafile cli list output: %w", err)
					}

					for _, library := range C.Libraries {
						if !slices.ContainsFunc(libraries, func(r SeafCliList) bool {
							return r.Id == library
						}) {
							return fmt.Errorf("library is missing from seafile client list output: %s", library)
						}
					}

					return nil
				}).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				return err
			}

			t.Log.Debugf(
				"Next repositories check in: %s",
				P.Health.StatusInterval.String(),
			)

			return nil
		})
}
