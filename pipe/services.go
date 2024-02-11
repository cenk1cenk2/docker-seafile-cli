package pipe

import (
	"path"
	"time"

	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func Services(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("services", "parent").
		SetJobWrapper(func(_ Job, _ *Task[Pipe]) Job {
			return tl.JobSequence(RunSeafDaemon(tl).Job(), tl.JobDelay(RunSeafileClient(tl).Job(), time.Second*3))
		})
}

func RunSeafDaemon(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("seaf-daemon").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				SEAFILE_CLI_EXE,
				"start",
				"-c",
				path.Join(t.Pipe.Seafile.DataLocation, "ccnet"),
			).
				EnableTerminator().
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				return err
			}

			t.Log.Infoln("Started Seafile Daemon.")

			return nil
		})
}

func RunSeafileClient(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("seafile-client").
		Set(func(t *Task[Pipe]) error {
			for _, library := range t.Pipe.Ctx.Libraries {
				func(library string) {
					t.CreateSubtask(library).
						Set(func(t *Task[Pipe]) error {
							// desync first
							t.CreateCommand(
								SEAFILE_CLI_EXE,
								"desync",
								"-d",
								path.Join(
									t.Pipe.Seafile.MountLocation,
									library,
								),
								"-c",
								path.Join(t.Pipe.Seafile.DataLocation, "ccnet"),
							).
								SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
								AddSelfToTheTask()

								// sync
							t.CreateCommand(
								SEAFILE_CLI_EXE,
								"sync",
								"-s",
								t.Pipe.Server.Url,
								"-u",
								t.Pipe.Credentials.Username,
								"-p",
								t.Pipe.Credentials.Password,
								"-l",
								library,
								"-d",
								path.Join(
									t.Pipe.Seafile.MountLocation,
									library,
								),
								"-c",
								path.Join(t.Pipe.Seafile.DataLocation, "ccnet"),
							).
								SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
								EnableTerminator().
								AddSelfToTheTask()

							return nil
						}).
						ShouldRunAfter(func(t *Task[Pipe]) error {
							return t.RunCommandJobAsJobSequence()
						}).
						AddSelfToTheParentAsParallel()
				}(library)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunSubtasks(); err != nil {
				return err
			}

			t.Log.Infoln("Started Seafile Client for library.")

			return nil
		})
}
