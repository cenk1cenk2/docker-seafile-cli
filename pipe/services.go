package pipe

import (
	"path"

	. "github.com/cenk1cenk2/plumber/v6"
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
		Set(func(t *Task) error {
			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"start",
				"-c",
				path.Join(P.Seafile.DataLocation, "ccnet"),
			).
				EnableTerminator().
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				return err
			}

			t.Log.Infoln("Started Seafile Daemon.")

			return nil
		})
}

func RunSeafileClient(tl *TaskList) *Task {
	return tl.CreateTask("seafile-client").
		Set(func(t *Task) error {
			for _, library := range C.Libraries {
				t.CreateSubtask(library).
					Set(func(t *Task) error {
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
							SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
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
							Set(func(c *Command) error {
								if P.Credentials.Token != "" {
									c.AppendArgs("-T", P.Credentials.Token)
								} else if P.Credentials.Password != "" {
									c.AppendArgs("-p", P.Credentials.Password)
								}

								return nil
							}).
							SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
							EnableTerminator().
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(t *Task) error {
						return t.RunCommandJobAsJobSequence()
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			if err := t.RunSubtasks(); err != nil {
				return err
			}

			t.Log.Infoln("Started Seafile Client for library.")

			return nil
		})
}
