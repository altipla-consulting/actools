
FROM golang:1.10
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

RUN apt-get update && \
    apt-get install -y unzip

RUN curl -sL -o protoc.zip https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip && \
    unzip protoc.zip -d /opt/protobuf && \
    rm protoc.zip && \
    chmod -R 0777 /opt/protobuf

RUN curl -sL -o master.zip https://github.com/googleapis/googleapis/archive/master.zip && \
    unzip master.zip -d /opt/googleapis && \
    rm master.zip && \
    chmod -R 0777 /opt/googleapis

ENV CACHE_BUST 1
ENV GOPATH /gotools
RUN go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u google.golang.org/grpc
ENV GOPATH /go

ENV PATH $PATH:/opt/protobuf/bin:/gotools/bin
WORKDIR /workspace
