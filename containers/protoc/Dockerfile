
FROM golang:1.8
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

RUN apt-get update && \
    apt-get install -y unzip

RUN curl -sL -o protoc.zip https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip && \
    unzip protoc.zip -d /opt/protobuf && \
    rm protoc.zip && \
    chmod -R 0777 /opt/protobuf
ENV PATH $PATH:/opt/protobuf/bin

ENV CACHE_BUST 1
RUN go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u google.golang.org/grpc

WORKDIR /workspace
