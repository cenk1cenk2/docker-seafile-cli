package pipe

import (
	"os"
	"path"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func Tasks(tl *TaskList) *Task {
	return tl.CreateTask("tasks", "parent").
		SetJobWrapper(func(_ Job, _ *Task) Job {
			return JobSequence(
				Secrets(tl).Job(),
				InitSeafile(tl).Job(),
				Setup(tl).Job(),
			)
		})
}

func Secrets(tl *TaskList) *Task {
	return tl.CreateTask("secrets").
		Set(func(t *Task) error {
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
		Set(func(t *Task) error {
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
					ShouldRunAfter(func(c *Command) error {
						c.Log.Infoln("Seafile data directory was empty so Seafile has been initiated.")

						return nil
					}).
					SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				t.Log.Debugln(err.Error())
			}

			return nil
		})
}

func Setup(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(t *Task) error {
			files, err := os.ReadDir(P.Seafile.MountLocation)

			if err != nil {
				return err
			}

			for _, file := range files {
				if file.IsDir() {
					C.Libraries = append(C.Libraries, file.Name())
				}
			}

			t.Log.Infof("Discovered libraries: %s", strings.Join(C.Libraries, ", "))

			return nil
		})
}
