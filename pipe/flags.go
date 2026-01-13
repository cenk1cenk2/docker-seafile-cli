package pipe

import (
	"time"

	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	category_health      = "Health"
	category_server      = "Server"
	category_credentials = "Credentials"
	category_seafile     = "Seafile"
)

var Flags = []cli.Flag{
	// category_health

	&cli.DurationFlag{
		Category: category_health,
		Name:     "health.check-interval",
		Usage:    "Health check interval for processes.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HEALTH_CHECK_INTERVAL"),
		),
		Value:       10 * time.Minute,
		Destination: &P.Health.CheckInterval,
	},

	&cli.DurationFlag{
		Category: category_health,
		Name:     "health.status-interval",
		Usage:    "Interval for outputting current status.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HEALTH_STATUS_INTERVAL"),
		),
		Value:       time.Hour,
		Destination: &P.Health.StatusInterval,
	},

	// category_server

	&cli.StringFlag{
		Category: category_server,
		Name:     "server.url",
		Usage:    "External url of the given Seafile server.",
		Required: true,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SERVER_URL"),
		),
		Destination: &P.Server.Url,
	},

	// category_credentials

	&cli.StringFlag{
		Category: category_credentials,
		Name:     "credentials.username",
		Usage:    "Email address of the user that owns the libraries.",
		Required: true,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("USERNAME"),
		),
		Destination: &P.Credentials.Username,
	},

	&cli.StringFlag{
		Category: category_credentials,
		Name:     "credentials.password",
		Usage:    "Password of the user that owns the libraries.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PASSWORD"),
		),
		Destination: &P.Credentials.Password,
	},

	&cli.StringFlag{
		Category: category_credentials,
		Name:     "credentials.token",
		Usage:    "Token of the user that owns the libraries.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TOKEN"),
		),
		Destination: &P.Credentials.Token,
	},

	// category_libraries

	&cli.StringFlag{
		Category: category_seafile,
		Name:     "seafile.mount-location",
		Usage:    "Mount location for the libraries.",
		Required: false,
		Value:    "/data",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SEAFILE_MOUNT_LOCATION"),
		),
		Destination: &P.Seafile.MountLocation,
	},

	&cli.StringFlag{
		Category: category_seafile,
		Name:     "seafile.data-location",
		Usage:    "Mount location for the data.",
		Required: false,
		Value:    "/seafile",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SEAFILE_DATA_LOCATION"),
		),
		Destination: &P.Seafile.DataLocation,
	},
}
