
FROM golang:1.14

RUN apt-get update && \
    apt-get install -y zip unzip

RUN mkdir /home/container && \
    chmod 0777 /home/container && \
    mkdir /data && \
    chmod 0777 /data

ENV PATH $PATH:/gotools/bin

CMD ["bash", "-c", "$APP"]
