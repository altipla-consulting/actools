
FROM golang:1.11

RUN apt-get update && \
    apt-get install -y zip unzip

RUN mkdir /home/container && \
   chmod 0777 /home/container

ENV CACHE_BUST 1
RUN go get -u github.com/cortesi/modd/cmd/modd

COPY run.sh /opt/run.sh

ENV PATH $PATH:/gotools/bin

CMD ["/opt/run.sh"]
