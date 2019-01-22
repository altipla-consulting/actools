
FROM node:10
MAINTAINER Ernesto Alejo <ernesto@altiplaconsulting.com>

RUN mkdir -p /home/container && \
    chmod 0777 /home/container

RUN npm install -g @vue/cli

COPY run.sh /opt/run.sh

ENV FORCE_COLOR true
CMD ["/opt/run.sh"]
