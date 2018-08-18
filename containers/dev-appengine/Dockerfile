
FROM google/cloud-sdk:latest
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

WORKDIR /workspace
CMD ["dev_appserver.py", "--datastore_path", "/home/local/datastore", "--host", "0.0.0.0", "--admin_host", "0.0.0.0", "--port", "8080", "app.yaml"]
