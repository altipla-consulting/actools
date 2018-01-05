
FROM golang:1.8
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

ENV CACHE_BUST 1
ENV GOPATH /gotools
RUN go get -u github.com/skelterjohn/rerun
ENV GOPATH /go

ENV PATH $PATH:/gotools/bin
WORKDIR /workspace
CMD ["bash"]
