
FROM golang:1.12

RUN apt-get update && \
    apt-get install -y unzip build-essential zlib1g-dev autoconf libtool

RUN curl -sL -o protoc.zip https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip && \
    unzip protoc.zip -d /opt/protobuf && \
    rm protoc.zip && \
    chmod -R 0777 /opt/protobuf

RUN curl -sL -o master.zip https://github.com/googleapis/googleapis/archive/master.zip && \
    unzip master.zip -d /opt/googleapis && \
    rm master.zip && \
    mv /opt/googleapis/googleapis-master/* /opt/googleapis && \
    rm -rf /opt/googleapis/googleapis-master && \
    chmod -R 0777 /opt/googleapis

WORKDIR /opt
RUN git clone -b $(curl -L https://grpc.io/release) https://github.com/grpc/grpc && \
   cd grpc && \
   git submodule update --init && \
   make grpc_php_plugin

ENV CACHE_BUST 6
ENV GOPATH /gotools
RUN go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u google.golang.org/grpc && \
    go get -u github.com/altipla-consulting/protoc-gen-grpc_browser
RUN cd /gotools/src/github.com/golang/protobuf && \
    git checkout v1.3.1 && \
    go install ./...
ENV GOPATH /go

ENV PATH $PATH:/opt/protobuf/bin:/gotools/bin
WORKDIR /workspace
