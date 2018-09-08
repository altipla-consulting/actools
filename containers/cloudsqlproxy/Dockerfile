
FROM launcher.gcr.io/google/debian9:latest

RUN apt-get update && \
    apt-get install -y wget

RUN wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 -O /opt/cloudsqlproxy && \
    chmod +x /opt/cloudsqlproxy

ENTRYPOINT ["/opt/cloudsqlproxy"]
