# docker-seafile-cli

Initiates the Seafile Client Daemon that will run in this container.

`docker-seafile-cli [FLAGS]`

## Flags

### CLI

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$LOG_LEVEL` | Define the log level for the application.  | `String`<br/>enum(&#34;PANIC&#34;, &#34;FATAL&#34;, &#34;WARNING&#34;, &#34;INFO&#34;, &#34;DEBUG&#34;, &#34;TRACE&#34;) | `false` | &#34;info&#34; |

### Credentials

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$USERNAME` | Email address of the user that owns the libraries. | `String` | `true` |  |
| `$PASSWORD` | Password of the user that owns the libraries. | `String` | `true` |  |

### Health

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$HEALTH_CHECK_INTERVAL` | Health check interval for processes. | `Duration` | `false` | 5m0s |
| `$HEALTH_STATUS_INTERVAL` | Interval for outputting current status. | `Duration` | `false` | 1h0m0s |

### Seafile

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$SEAFILE_MOUNT_LOCATION` | Mount location for the libraries. | `String` | `false` | &#34;/data&#34; |
| `$SEAFILE_DATA_LOCATION` | Mount location for the data. | `String` | `false` | &#34;/seafile&#34; |

### Server

| Flag / Environment |  Description   |  Type    | Required | Default |
|---------------- | --------------- | --------------- |  --------------- |  --------------- |
| `$SERVER_URL` | External url of the given Seafile server. | `String` | `true` |  |
