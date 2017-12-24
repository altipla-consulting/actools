
FROM golang:1.8
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

RUN apt-get update && \
    apt-get install -y curl

RUN curl -L -S -o glide.tar.gz https://github.com/Masterminds/glide/releases/download/v0.13.1/glide-v0.13.1-linux-amd64.tar.gz && \
    tar -xzf glide.tar.gz && \
    mv linux-amd64/glide /usr/bin/glide && \
    rm -r linux-amd64

ENV GLIDE_HOME /tmp

WORKDIR /staging
