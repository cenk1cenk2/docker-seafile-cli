package pipe

import (
	"fmt"
	"os"
	"strings"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func Tasks(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("tasks", "parent").
		SetJobWrapper(func(job Job) Job {
			return TL.JobSequence(
				Setup(tl).Job(),
			)
		})
}

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		Set(func(t *Task[Pipe]) error {
			files, err := os.ReadDir(t.Pipe.Libraries.MountLocation)

			if err != nil {
				return err
			}

			for _, file := range files {
				if file.IsDir() {
					t.Pipe.Ctx.Libraries = append(t.Pipe.Ctx.Libraries, file.Name())
				}
			}

			if len(t.Pipe.Ctx.Libraries) == 0 {
				return fmt.Errorf("Please mount your libraries as folders to: %s", t.Pipe.Libraries.MountLocation)
			}

			t.Log.Infof("Discovered libraries: %s", strings.Join(t.Pipe.Ctx.Libraries, ", "))

			return nil
		})
}
