
FROM golang:1.17.3

RUN apt-get update && \
    apt-get install -y unzip build-essential zlib1g-dev autoconf libtool cmake

RUN curl -sL -o protoc.zip https://github.com/google/protobuf/releases/download/v3.19.1/protoc-3.19.1-linux-x86_64.zip && \
    unzip protoc.zip -d /opt/protobuf && \
    rm protoc.zip && \
    chmod -R 0777 /opt/protobuf

RUN curl -sL -o master.zip https://github.com/googleapis/googleapis/archive/master.zip && \
    unzip master.zip -d /opt/googleapis && \
    rm master.zip && \
    mv /opt/googleapis/googleapis-master/* /opt/googleapis && \
    rm -rf /opt/googleapis/googleapis-master && \
    chmod -R 0777 /opt/googleapis

RUN curl -L https://github.com/altipla-consulting/protoc-gen-grpc_browser/releases/download/v0.4.2/protoc-gen-grpc_browser_v0.4.2_linux_amd64 -o /usr/bin/protoc-gen-grpc_browser && \
    chmod +x /usr/bin/protoc-gen-grpc_browser

ENV CACHE_BUST 8
ENV GOPATH /gotools
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.7.2 && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.7.2
ENV GOPATH /go

ENV PATH $PATH:/opt/protobuf/bin:/gotools/bin
WORKDIR /workspace
