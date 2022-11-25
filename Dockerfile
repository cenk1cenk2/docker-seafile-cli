# syntax=docker/dockerfile-upstream:master-labs
FROM debian:bullseye-slim

ARG BUILDOS
ARG BUILDARCH

RUN \
  apt-get update && \
  apt-get install gnupg2 wget -y && \
  # install tini
  apt-get install tini -y && \
  # install seafile client
  wget https://linux-clients.seafile.com/seafile.asc -O /usr/share/keyrings/seafile-keyring.asc && \
  echo 'deb [arch=amd64 signed-by=/usr/share/keyrings/seafile-keyring.asc] https://linux-clients.seafile.com/seafile-deb/bullseye/ stable main' | tee /etc/apt/sources.list.d/seafile.list && \
  apt-get install -y seafile-cli procps curl grep && \
  rm -rf /var/lib/apt/lists/*

COPY --chmod=777 ./dist/pipe-${BUILDOS}-${BUILDARCH} /usr/bin/pipe

RUN \
  # smoke test
  pipe --help

WORKDIR /data

ENTRYPOINT [ "tini", "pipe" ]
