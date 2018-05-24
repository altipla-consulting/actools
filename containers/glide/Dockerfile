
FROM golang:1.10
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

RUN apt-get update && \
    apt-get install -y curl

RUN curl -L -S -o glide.tar.gz https://github.com/Masterminds/glide/releases/download/v0.13.1/glide-v0.13.1-linux-amd64.tar.gz && \
    tar -xzf glide.tar.gz && \
    mv linux-amd64/glide /usr/bin/glide && \
    rm -r linux-amd64

RUN groupadd --gid 1000 -r localgrp -o

RUN useradd --system --uid=1000 --gid=1000 --home-dir /home/local1000 local1000 -o && \
    mkdir /home/local1000 && \
    chown local1000:localgrp /home/local1000
COPY gitconfig /home/local1000/.gitconfig
COPY config /home/local1000/.ssh/config

RUN useradd --system --uid=1001 --gid=1000 --home-dir /home/local1001 local1001 -o && \
    mkdir /home/local1001 && \
    chown local1001:localgrp /home/local1001
COPY gitconfig /home/local1001/.gitconfig
COPY config /home/local1001/.ssh/config

ENV GLIDE_HOME /tmp

WORKDIR /workspace
