package pipe

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
)

func Tasks(tl *TaskList) *Task {
	return tl.CreateTask("tasks", "parent").
		SetJobWrapper(func(_ Job, _ *Task) Job {
			return JobSequence(
				Secrets(tl).Job(),
				InitSeafile(tl).Job(),
				ConfigureSeafile(tl).Job(),
				Setup(tl).Job(),
			)
		})
}

func Secrets(tl *TaskList) *Task {
	return tl.CreateTask("secrets").
		Set(func(_ context.Context, t *Task) error {
			t.Plumber.AppendSecrets(P.Credentials.Username)

			if P.Credentials.Password != "" {
				t.Plumber.AppendSecrets(P.Credentials.Password)
			}
			if P.Credentials.Token != "" {
				t.Plumber.AppendSecrets(P.Credentials.Token)
			}

			return nil
		})
}

func InitSeafile(tl *TaskList) *Task {
	return tl.CreateTask("seafile").
		Set(func(_ context.Context, t *Task) error {
			files, err := os.ReadDir(P.Seafile.DataLocation)

			if err != nil {
				return err
			}

			if len(files) == 0 {
				t.CreateCommand(
					SEAFILE_CLI_EXE,
					"init",
					"-d",
					P.Seafile.DataLocation,
					"-c",
					path.Join(P.Seafile.DataLocation, "ccnet"),
				).
					ShouldRunAfter(func(_ context.Context, c *Command) error {
						c.Log.Info("Seafile data directory was empty so Seafile has been initiated.")

						return nil
					}).
					SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			if err := t.RunCommandJobAsJobSequence(ctx); err != nil {
				t.Log.Debug(err.Error())
			}

			return nil
		})
}

func ConfigureSeafile(tl *TaskList) *Task {
	return tl.CreateTask("seafile", "config").
		Set(func(_ context.Context, t *Task) error {
			if P.Seafile.Umask == "" {
				return nil
			}

			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"config",
				"-c",
				path.Join(P.Seafile.DataLocation, "ccnet"),
				"-k",
				"umask",
				"-v",
				P.Seafile.Umask,
			).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					c.Log.Info(fmt.Sprintf("Configured Seafile umask: %s", P.Seafile.Umask))

					return nil
				}).
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func Setup(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(_ context.Context, t *Task) error {
			files, err := os.ReadDir(P.Seafile.MountLocation)

			if err != nil {
				return err
			}

			for _, file := range files {
				if file.IsDir() {
					C.Libraries = append(C.Libraries, file.Name())
				}
			}

			t.Log.Info(fmt.Sprintf("Discovered libraries: %s", strings.Join(C.Libraries, ", ")))

			return nil
		})
}
