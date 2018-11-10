
FROM golang:1.11

ENV CACHE_BUST 4
RUN go get -u github.com/altipla-consulting/migrator/cmd/init-migrator
RUN go get -u github.com/altipla-consulting/migrator/cmd/migrator

WORKDIR /workspace
