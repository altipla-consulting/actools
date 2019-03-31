
FROM debian:stretch
ARG CLOUD_SDK_VERSION=239.0.0
ENV CLOUD_SDK_VERSION=$CLOUD_SDK_VERSION
ENV PATH "$PATH:/opt/google-cloud-sdk/bin/"
RUN apt-get -y update && apt-get install -y \
        curl \
        gcc \
        python-dev \
        python-setuptools \
        apt-transport-https \
        lsb-release \
        openssh-client \
        git \
        gnupg \
        openjdk-8-jre \
    && easy_install -U pip && \
    pip install -U crcmod   && \
    mkdir -p /opt && cd /opt && \
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
    rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz
RUN gcloud components install -q cloud-firestore-emulator beta && \
    gcloud config set core/disable_usage_reporting true && \
    gcloud config set component_manager/disable_update_check true && \
    gcloud config set metrics/environment github_docker_image && \
    gcloud --version
VOLUME ["/root/.config"]

RUN mkdir /home/container

WORKDIR /workspace
CMD ["gcloud", "beta", "emulators", "firestore", "start", "--host-port", "127.0.0.1:8080"]
