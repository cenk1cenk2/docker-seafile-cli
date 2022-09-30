package main

import (
	"github.com/urfave/cli/v2"

	pipe "gitlab.kilic.dev/docker/seafile-cli/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func main() {
	p := Plumber{
		DocsFile:               "CLI.md",
		DocsExcludeFlags:       true,
		DocsExcludeHelpCommand: true,
	}

	p.New(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       p.AppendFlags(pipe.Flags),
				Before: func(ctx *cli.Context) error {
					p.EnableTerminator()

					return nil
				},
				Action: func(ctx *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.New(p).SetCliContext(ctx).Job(),
					)
				},
			}
		}).
		Run()
}