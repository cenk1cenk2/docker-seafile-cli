package pipe

import (
	"context"
	"path"

	. "github.com/cenk1cenk2/plumber/v7"
)

func Services(tl *TaskList) *Task {
	return tl.CreateTask("services", "parent").
		SetJobWrapper(func(_ Job, _ *Task) Job {
			return JobSequence(
				RunSeafDaemon(tl).Job(),
				RunSeafileClient(tl).Job(),
			)
		})
}

func RunSeafDaemon(tl *TaskList) *Task {
	return tl.CreateTask("seaf-daemon").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"start",
				"-c",
				path.Join(P.Seafile.DataLocation, "ccnet"),
			).
				EnableTerminator().
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			if err := t.RunCommandJobAsJobSequence(ctx); err != nil {
				return err
			}

			t.Log.Info("Started Seafile Daemon.")

			return nil
		})
}

func RunSeafileClient(tl *TaskList) *Task {
	return tl.CreateTask("seafile-client").
		Set(func(_ context.Context, t *Task) error {
			for _, library := range C.Libraries {
				t.CreateSubtask(library).
					Set(func(_ context.Context, t *Task) error {
						// desync first
						t.CreateCommand(
							SEAFILE_CLI_EXE,
							"desync",
							"-d",
							path.Join(
								P.Seafile.MountLocation,
								library,
							),
							"-c",
							path.Join(P.Seafile.DataLocation, "ccnet"),
						).
							SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
							AddSelfToTheTask()

							// sync
						t.CreateCommand(
							SEAFILE_CLI_EXE,
							"sync",
							"-s",
							P.Server.Url,
							"-l",
							library,
							"-d",
							path.Join(
								P.Seafile.MountLocation,
								library,
							),
							"-c",
							path.Join(P.Seafile.DataLocation, "ccnet"),
							"-u",
							P.Credentials.Username,
						).
							Set(func(_ context.Context, c *Command) error {
								if P.Credentials.Token != "" {
									c.AppendArgs("-T", P.Credentials.Token)
								} else if P.Credentials.Password != "" {
									c.AppendArgs("-p", P.Credentials.Password)
								}

								return nil
							}).
							SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
							EnableTerminator().
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(ctx context.Context, t *Task) error {
						return t.RunCommandJobAsJobSequence(ctx)
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			if err := t.RunSubtasks(ctx); err != nil {
				return err
			}

			t.Log.Info("Started Seafile Client for library.")

			return nil
		})
}
