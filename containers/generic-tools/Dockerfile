
FROM launcher.gcr.io/google/debian9:latest

RUN apt-get update && \
  apt-get install -y curl jq

ENV CACHE_BUST validate-http-proto-v0.1.3
RUN curl -sSL https://git.io/fjE3d | bash

ENV CACHE_BUST devmails-v0.1.0
RUN curl -sSL https://git.io/fjuM7 | bash

WORKDIR /workspace
