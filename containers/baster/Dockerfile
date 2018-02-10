
FROM golang:1.8
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

ENV CACHE_BUST 2
RUN go get -u github.com/altipla-consulting/baster/cmd/baster

WORKDIR /workspace
CMD ["baster"]
