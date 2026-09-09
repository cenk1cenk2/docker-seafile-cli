# docker-seafile-cli

Initiates the Seafile Client Daemon that will run in this container.

`docker-seafile-cli [FLAGS]`

## Flags

**CLI**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `"info"` |
| `$ENV_FILE` | Environment files to inject. | `string[]` |  |

**Credentials**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| **`$USERNAME`**\* | Email address of the user that owns the libraries. | `string` |  |
| `$PASSWORD` | Password of the user that owns the libraries. | `string` |  |
| `$TOKEN` | Token of the user that owns the libraries. | `string` |  |

\* required

**Health**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$HEALTH_STATUS_INTERVAL` | Interval for outputting current status. | `duration` | `5m0s` |
| `$HEALTH_EXIT_ON_FAILURE` | Exit the process when a health check fails, so the orchestrator can restart the container. | `bool` | `true` |

**Seafile**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| `$SEAFILE_MOUNT_LOCATION` | Mount location for the libraries. | `string` | `"/data"` |
| `$SEAFILE_DATA_LOCATION` | Mount location for the data. | `string` | `"/seafile"` |
| `$SEAFILE_UMASK` | Umask to configure for the Seafile client. | `string` |  |

**Server**

| Flag / Environment | Description | Type | Default |
| --- | --- | --- | --- |
| **`$SERVER_URL`**\* | External url of the given Seafile server. | `string` |  |

\* required
