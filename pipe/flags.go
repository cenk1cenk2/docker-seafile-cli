package pipe

import (
	"time"

	"github.com/urfave/cli/v2"
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
		Category:    category_health,
		Name:        "health.check-interval",
		Usage:       "Health check interval for processes.",
		Required:    false,
		EnvVars:     []string{"HEALTH_CHECK_INTERVAL"},
		Value:       5 * time.Minute,
		Destination: &TL.Pipe.Health.CheckInterval,
	},

	&cli.DurationFlag{
		Category:    category_health,
		Name:        "health.status-interval",
		Usage:       "Interval for outputting current status.",
		Required:    false,
		EnvVars:     []string{"HEALTH_STATUS_INTERVAL"},
		Value:       time.Hour,
		Destination: &TL.Pipe.Health.StatusInterval,
	},

	// category_server

	&cli.StringFlag{
		Category:    category_server,
		Name:        "server.url",
		Usage:       "External url of the given Seafile server.",
		Required:    true,
		EnvVars:     []string{"SERVER_URL"},
		Destination: &TL.Pipe.Server.Url,
	},

	// category_credentials

	&cli.StringFlag{
		Category:    category_credentials,
		Name:        "credentials.username",
		Usage:       "Email address of the user that owns the libraries.",
		Required:    true,
		EnvVars:     []string{"USERNAME"},
		Destination: &TL.Pipe.Credentials.Username,
	},

	&cli.StringFlag{
		Category:    category_credentials,
		Name:        "credentials.password",
		Usage:       "Password of the user that owns the libraries.",
		Required:    true,
		EnvVars:     []string{"PASSWORD"},
		Destination: &TL.Pipe.Credentials.Password,
	},

	// category_libraries

	&cli.StringFlag{
		Category:    category_seafile,
		Name:        "seafile.mount-location",
		Usage:       "Mount location for the libraries.",
		Required:    false,
		Value:       "/data",
		EnvVars:     []string{"SEAFILE_MOUNT_LOCATION"},
		Destination: &TL.Pipe.Seafile.MountLocation,
	},

	&cli.StringFlag{
		Category:    category_seafile,
		Name:        "seafile.data-location",
		Usage:       "Mount location for the data.",
		Required:    false,
		Value:       "/seafile",
		EnvVars:     []string{"SEAFILE_DATA_LOCATION"},
		Destination: &TL.Pipe.Seafile.DataLocation,
	},
}
