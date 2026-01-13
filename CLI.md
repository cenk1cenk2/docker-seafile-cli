# docker-seafile-cli

Initiates the Seafile Client Daemon that will run in this container.

`docker-seafile-cli [FLAGS]`

## Flags

**CLI**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application. | `string`<br/>`enum("panic", "fatal", "warn", "info", "debug", "trace")` | `false` | <code>"info"</code> |
| `$ENV_FILE` | Environment files to inject. | `string[]` | `false` | <code></code> |

**Credentials**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$USERNAME` | Email address of the user that owns the libraries. | `string` | `true` | <code></code> |
| `$PASSWORD` | Password of the user that owns the libraries. | `string` | `false` | <code></code> |
| `$TOKEN` | Token of the user that owns the libraries. | `string` | `false` | <code></code> |

**Health**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HEALTH_CHECK_INTERVAL` | Health check interval for processes. | `duration` | `false` | <code>10m0s</code> |
| `$HEALTH_STATUS_INTERVAL` | Interval for outputting current status. | `duration` | `false` | <code>1h0m0s</code> |

**Seafile**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$SEAFILE_MOUNT_LOCATION` | Mount location for the libraries. | `string` | `false` | <code>"/data"</code> |
| `$SEAFILE_DATA_LOCATION` | Mount location for the data. | `string` | `false` | <code>"/seafile"</code> |

**Server**

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$SERVER_URL` | External url of the given Seafile server. | `string` | `true` | <code></code> |
