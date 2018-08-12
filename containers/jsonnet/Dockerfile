
FROM golang:1.10

ENV CACHE_BUST 1
RUN go get -u github.com/google/go-jsonnet/jsonnet
RUN go get -u github.com/ksonnet/kubecfg

WORKDIR /workspace
