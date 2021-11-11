
FROM golang:1.17.3

ENV CACHE_BUST 5
RUN go get -u github.com/altipla-consulting/migrator/cmd/init-migrator
RUN go get -u github.com/altipla-consulting/migrator/cmd/migrator

WORKDIR /workspace
