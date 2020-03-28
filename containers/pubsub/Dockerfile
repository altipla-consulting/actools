
FROM google/cloud-sdk:latest

RUN mkdir /home/container

WORKDIR /workspace
CMD ["gcloud", "beta", "emulators", "pubsub", "start", "--project=local"]
