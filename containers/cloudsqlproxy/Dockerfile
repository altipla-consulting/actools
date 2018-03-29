
FROM google/debian:jessie
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

RUN apt-get update && \
    apt-get install -y wget

RUN wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 -O /opt/cloudsqlproxy && \
    chmod +x /opt/cloudsqlproxy

CMD ["/opt/cloudsqlproxy"]
