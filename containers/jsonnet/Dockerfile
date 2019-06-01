
FROM golang:1.11

RUN mkdir /home/container

ENV CACHE_BUST 1
RUN go get -u github.com/google/go-jsonnet/cmd/jsonnet
RUN go get -u github.com/ksonnet/kubecfg

WORKDIR /workspace
