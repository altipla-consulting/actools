
FROM golang:1.13

RUN apt-get update && \
    apt-get install -y jq

RUN mkdir /home/container && \
    chmod 0777 /home/container

WORKDIR /go/bin

RUN curl -L -o reloader $(curl --silent 'https://api.github.com/repos/altipla-consulting/reloader/releases/tags/v0.3.4' | jq -r '.assets[0].browser_download_url') && \
    chmod +x /go/bin/reloader
