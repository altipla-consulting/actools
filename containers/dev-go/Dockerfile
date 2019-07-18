
FROM golang:1.12

RUN apt-get update && \
    apt-get install -y zip unzip

RUN mkdir /home/container && \
    chmod 0777 /home/container && \
    mkdir /data && \
    chmod 0777 /data

ENV CACHE_BUST 1
ENV GOPATH /gotools
RUN go get -u github.com/cortesi/modd/cmd/modd
ENV GOPATH /go

COPY run.sh /opt/run.sh

ENV PATH $PATH:/gotools/bin

CMD ["/opt/run.sh"]
