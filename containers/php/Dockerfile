
FROM php:7

RUN apt-get update && \
    apt-get install -y zlib1g-dev wget git unzip

RUN pecl install grpc && \
    docker-php-ext-enable grpc

WORKDIR /usr/bin
COPY install.sh /opt/install.sh
RUN /opt/install.sh

RUN wget -O phpunit https://phar.phpunit.de/phpunit-7.phar && \
    chmod +x phpunit
