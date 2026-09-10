package pipe

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"slices"
	"strings"
	"time"

	. "github.com/cenk1cenk2/plumber/v7"
)

func HealthCheck(tl *TaskList) *Task {
	return tl.CreateTask("health", "parent").
		SetJobWrapper(func(_ Job, t *Task) Job {
			check := JobParallel(
				HealthCheckStatus(tl).Job(),
				HealthCheckRepositories(tl).Job(),
			)

			if P.Health.ExitOnFailure {
				run := check

				check = func(ctx context.Context) error {
					if err := run(ctx); err != nil {
						t.SendFatal(err)
					}

					return nil
				}
			} else {
				check = GuardResume(check, t.Log)
			}

			return JobBackground(
				JobDelay(
					JobLoopWithWaitAfter(
						check,
						P.Health.StatusInterval,
					),
					15*time.Second,
				),
				t.Log,
			)
		})
}

func HealthCheckStatus(tl *TaskList) *Task {
	return tl.CreateTask("health", "status").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"status",
				"-c",
				path.Join(P.Seafile.DataLocation, "ccnet"),
			).
				SetLogLevel(LogLevelInfo, LogLevelWarn, LogLevelDebug).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			if err := t.RunCommandJobAsJobSequence(ctx); err != nil {
				return err
			}

			t.Log.Debug(fmt.Sprintf(
				"Next status check in: %s",
				P.Health.StatusInterval.String(),
			))

			return nil
		})
}

func HealthCheckRepositories(tl *TaskList) *Task {
	return tl.CreateTask("health", "repositories").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"list",
				"-c",
				path.Join(P.Seafile.DataLocation, "ccnet"),
				"--json",
			).
				EnableStreamRecording().
				ShouldRunAfter(func(_ context.Context, c *Command) error {
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
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			if err := t.RunCommandJobAsJobSequence(ctx); err != nil {
				return err
			}

			t.Log.Debug(fmt.Sprintf(
				"Next repositories check in: %s",
				P.Health.StatusInterval.String(),
			))

			return nil
		})
}
