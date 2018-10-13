
FROM google/cloud-sdk:220.0.0

RUN apt-get update && \
    apt-get install -y wget

RUN cd /tmp && \
    wget -q https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz && \
    tar -xvf go1.11.1.linux-amd64.tar.gz && \
    mv go /usr/local && \
    rm /tmp/go1.11.1.linux-amd64.tar.gz

RUN mkdir -p /home/container && \
    chmod 0777 /home/container

COPY run.sh /opt/run.sh

ENV PATH /usr/local/go/bin:$PATH

CMD ["/opt/run.sh"]
