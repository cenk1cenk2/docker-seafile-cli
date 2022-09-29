package pipe

import (
	"fmt"
	"path"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func Services(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("services", "parent").
		SetJobWrapper(func(job Job) Job {
			return TL.JobSequence(RunSeafDaemon(tl).Job(), RunSeafileClient(tl).Job())
		})
}

func RunSeafDaemon(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("seaf-daemon").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand("seaf-daemon").
				EnableTerminator().
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
							t.CreateCommand(
								SEAFILE_CLI_EXE,
								"sync",
								"-l",
								library,
								"-d",
								path.Join(
									t.Pipe.Libraries.MountLocation,
									library,
								),
								"-u",
								fmt.Sprintf("'%s'", t.Pipe.Credentials.Username),
								"-p",
								fmt.Sprintf("'%s'", t.Pipe.Credentials.Password),
							).
								EnableTerminator().
								AddSelfToTheTask()

							return nil
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
