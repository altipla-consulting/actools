
FROM golang:1.13

RUN apt-get update && \
    apt-get install -y unzip build-essential zlib1g-dev autoconf libtool

RUN curl -sL -o protoc.zip https://github.com/google/protobuf/releases/download/v3.11.4/protoc-3.11.4-linux-x86_64.zip && \
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

RUN curl -L https://github.com/altipla-consulting/protoc-gen-grpc_browser/releases/download/v0.4.2/protoc-gen-grpc_browser_v0.4.2_linux_amd64 -o /usr/bin/protoc-gen-grpc_browser && \
    chmod +x /usr/bin/protoc-gen-grpc_browser

ENV CACHE_BUST 8
ENV GOPATH /gotools
RUN go get github.com/golang/protobuf/proto && \
    go get github.com/golang/protobuf/protoc-gen-go && \
    go get google.golang.org/grpc && \
    go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
    go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN cd /gotools/src/github.com/golang/protobuf && \
    git checkout v1.3.4 && \
    go install ./...
ENV GOPATH /go

ENV PATH $PATH:/opt/protobuf/bin:/gotools/bin
WORKDIR /workspace
