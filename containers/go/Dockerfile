
FROM golang:1.14.7

RUN apt-get update && \
    apt-get install -y zip unzip jq

RUN mkdir /home/container && \
    chmod 0777 /home/container && \
    mkdir /data && \
    chmod 0777 /data

RUN groupadd --gid 1000 -r localgrp -o && \
    useradd --system --uid=1000 --gid=1000 --home-dir /home/container local1000 -o && \
    useradd --system --uid=1001 --gid=1000 --home-dir /home/container local1001 -o

RUN curl https://tools.altipla.consulting/tools/reloader > /usr/bin/reloader && \
    chmod +x /usr/bin/reloader

COPY gitconfig /home/container/.gitconfig
COPY gitconfig /root/.gitconfig

COPY sshconfig /home/container/.ssh/config
COPY sshconfig /root/.ssh/config

ENV HOME /home/container

WORKDIR /workspace
