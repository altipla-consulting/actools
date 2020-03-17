
FROM golang:1.13

RUN apt-get update && \
    apt-get install -y zip unzip

RUN mkdir /home/container && \
    chmod 0777 /home/container && \
    mkdir /data && \
    chmod 0777 /data

RUN groupadd --gid 1000 -r localgrp -o && \
    useradd --system --uid=1000 --gid=1000 --home-dir /home/container local1000 -o && \
    useradd --system --uid=1001 --gid=1000 --home-dir /home/container local1001 -o

ENV CACHE_BUST v0.3.1
RUN go get github.com/altipla-consulting/reloader && \
    mv /go/bin/reloader /usr/local/bin/reloader

COPY gitconfig /home/container/.gitconfig
COPY gitconfig /root/.gitconfig

COPY sshconfig /home/container/.ssh/config
COPY sshconfig /root/.ssh/config

ENV HOME /home/container

WORKDIR /workspace
