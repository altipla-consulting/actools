
FROM golang:1.14

RUN mkdir /home/container

ENV CACHE_BUST 1
RUN go get -u github.com/google/go-jsonnet/cmd/jsonnet
RUN go get -u github.com/bitnami/kubecfg

WORKDIR /workspace
