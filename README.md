# cenk1cenk2/seafile-cli

[![pipeline status](https://gitlab.kilic.dev/docker/seafile-cli/badges/master/pipeline.svg)](https://gitlab.kilic.dev/docker/seafile-cli/-/commits/master) [![Docker Pulls](https://img.shields.io/docker/pulls/cenk1cenk2/seafile-cli)](https://hub.docker.com/repository/docker/cenk1cenk2/seafile-cli) [![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/cenk1cenk2/seafile-cli)](https://hub.docker.com/repository/docker/cenk1cenk2/seafile-cli) [![Docker Image Version (latest by date)](https://img.shields.io/docker/v/cenk1cenk2/seafile-cli)](https://hub.docker.com/repository/docker/cenk1cenk2/seafile-cli) [![GitHub last commit](https://img.shields.io/github/last-commit/cenk1cenk2/seafile-cli)](https://github.com/cenk1cenk2/seafile-cli)

## Description

Seafile client inside a Docker container that can synchronize multiple libraries for the given user credentials.

---

- [CLI Documentation](./CLI.md)
<!-- toc -->

<!-- tocstop -->

---

## Setup

Mount your libraries to `$SEAFILE_MOUNT_LOCATION` in subfolders with the library UUIDs that can be obtained through `https://seafile.yourdomain.com/libraries/$UUID`. It will automatically synchronize the libraries under the given folder for the given user credentials.

## Environment Variables

### CLI

| Flag / Environment | Description                               | Type                                                                       | Required | Default |
| ------------------ | ----------------------------------------- | -------------------------------------------------------------------------- | -------- | ------- |
| `$LOG_LEVEL`       | Define the log level for the application. | `String`<br/>`enum("PANIC", "FATAL", "WARNING", "INFO", "DEBUG", "TRACE")` | `false`  | info    |
| `$ENV_FILE`        | Environment files to inject.              | `StringSlice`                                                              | `false`  |         |

### Credentials

| Flag / Environment | Description                                        | Type     | Required | Default |
| ------------------ | -------------------------------------------------- | -------- | -------- | ------- |
| `$USERNAME`        | Email address of the user that owns the libraries. | `String` | `true`   |         |
| `$PASSWORD`        | Password of the user that owns the libraries.      | `String` | `true`   |         |

### Health

| Flag / Environment        | Description                             | Type       | Required | Default |
| ------------------------- | --------------------------------------- | ---------- | -------- | ------- |
| `$HEALTH_CHECK_INTERVAL`  | Health check interval for processes.    | `Duration` | `false`  | 10m     |
| `$HEALTH_STATUS_INTERVAL` | Interval for outputting current status. | `Duration` | `false`  | 1h      |

### Seafile

| Flag / Environment        | Description                       | Type     | Required | Default  |
| ------------------------- | --------------------------------- | -------- | -------- | -------- |
| `$SEAFILE_MOUNT_LOCATION` | Mount location for the libraries. | `String` | `false`  | /data    |
| `$SEAFILE_DATA_LOCATION`  | Mount location for the data.      | `String` | `false`  | /seafile |

### Server

| Flag / Environment | Description                               | Type     | Required | Default |
| ------------------ | ----------------------------------------- | -------- | -------- | ------- |
| `$SERVER_URL`      | External url of the given Seafile server. | `String` | `true`   |         |

## Deploy

You can check out the example setup inside `docker-compose.yml` to see how this container can become operational. Please mind the environment variables for the configuration as well as the section about the configuration files and their generation.
