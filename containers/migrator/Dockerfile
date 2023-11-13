
FROM golang:1.21.4
ENV CACHE_BUST 5
RUN go install github.com/altipla-consulting/migrator/cmd/init-migrator@latest
RUN go install github.com/altipla-consulting/migrator/cmd/migrator@latest

WORKDIR /workspace
