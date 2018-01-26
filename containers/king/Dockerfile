
FROM golang:1.8
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

RUN apt-get update && \
    apt-get install -y unzip

RUN curl -sL -o protoc.zip https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip && \
    unzip protoc.zip -d /opt/protobuf && \
    rm protoc.zip && \
    chmod -R 0777 /opt/protobuf

ENV CACHE_BUST 2
ENV GOPATH /gotools
RUN go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u github.com/altipla-consulting/king/protoc-gen-king_browser && \
    go get -u github.com/altipla-consulting/king/protoc-gen-king_go
ENV GOPATH /go

RUN ln -s /opt/protobuf/bin/protoc /usr/bin/king

ENV PATH $PATH:/bin:/gotools/bin
WORKDIR /workspace
