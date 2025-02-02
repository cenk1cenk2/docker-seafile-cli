package pipe

import (
	"os"
	"path"
	"strings"

	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func Tasks(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tasks", "parent").
		SetJobWrapper(func(_ Job, _ *Task[Pipe]) Job {
			return tl.JobSequence(
				Secrets(tl).Job(),
				InitSeafile(tl).Job(),
				Setup(tl).Job(),
			)
		})
}

func Secrets(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("secrets").
		Set(func(t *Task[Pipe]) error {
			t.Plumber.AppendSecrets(t.Pipe.Credentials.Username)

			if t.Pipe.Credentials.Password != "" {
				t.Plumber.AppendSecrets(t.Pipe.Credentials.Password)
			}
			if t.Pipe.Credentials.Token != "" {
				t.Plumber.AppendSecrets(t.Pipe.Credentials.Token)
			}

			return nil
		})
}

func InitSeafile(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("seafile").
		Set(func(t *Task[Pipe]) error {
			files, err := os.ReadDir(t.Pipe.Seafile.DataLocation)

			if err != nil {
				return err
			}

			if len(files) == 0 {
				t.CreateCommand(
					SEAFILE_CLI_EXE,
					"init",
					"-d",
					t.Pipe.Seafile.DataLocation,
					"-c",
					path.Join(t.Pipe.Seafile.DataLocation, "ccnet"),
				).
					ShouldRunAfter(func(c *Command[Pipe]) error {
						c.Log.Infoln("Seafile data directory was empty so Seafile has been initiated.")

						return nil
					}).
					SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				t.Log.Debugln(err.Error())
			}

			return nil
		})
}

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		Set(func(t *Task[Pipe]) error {
			files, err := os.ReadDir(t.Pipe.Seafile.MountLocation)

			if err != nil {
				return err
			}

			for _, file := range files {
				if file.IsDir() {
					t.Pipe.Ctx.Libraries = append(t.Pipe.Ctx.Libraries, file.Name())
				}
			}

			t.Log.Infof("Discovered libraries: %s", strings.Join(t.Pipe.Ctx.Libraries, ", "))

			return nil
		})
}
