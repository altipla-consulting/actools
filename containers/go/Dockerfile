
FROM golang:1.14

RUN apt-get update && \
    apt-get install -y zip unzip

RUN groupadd --gid 1000 -r localgrp -o && \
    useradd --system --uid=1000 --gid=1000 --home-dir /home/container local1000 -o && \
    useradd --system --uid=1001 --gid=1000 --home-dir /home/container local1001 -o

COPY gitconfig /home/container/.gitconfig
COPY gitconfig /root/.gitconfig

COPY sshconfig /home/container/.ssh/config
COPY sshconfig /root/.ssh/config

WORKDIR /workspace
