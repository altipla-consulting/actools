
FROM launcher.gcr.io/google/debian9:latest

RUN apt-get update && \
  apt-get install -y curl jq

# validate-http-proto
RUN curl -sSL https://git.io/fjE3d | bash

# devmails
RUN curl -sSL https://git.io/fjuM7 | bash

WORKDIR /workspace
