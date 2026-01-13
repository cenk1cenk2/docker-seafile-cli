# syntax=docker/dockerfile:1.20
FROM debian:trixie-slim

ARG BUILDOS
ARG BUILDARCH

RUN \
  apt-get update && \
  apt-get install --no-install-recommends -y \
  ca-certificates \
  tini \
  seafile-cli \
  procps \
  curl \
  grep && \
  update-ca-certificates && \
  rm -rf /var/lib/apt/lists/*

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help

WORKDIR /data

ENTRYPOINT [ "tini", "pipe" ]
