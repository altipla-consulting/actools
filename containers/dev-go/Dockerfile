
FROM golang:1.10
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

ENV CACHE_BUST 1
ENV GOPATH /gotools
RUN go get -u github.com/cortesi/modd/cmd/modd
ENV GOPATH /go

COPY run.sh /opt/run.sh

ENV PATH $PATH:/gotools/bin

CMD ["/opt/run.sh"]
