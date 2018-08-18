
FROM google/cloud-sdk:latest

RUN mkdir /opt/local && \
    chmod 0777 /opt/local

COPY run.sh /opt/run.sh

ENV GOPATH /go

WORKDIR /workspace
CMD ["/opt/run.sh"]
