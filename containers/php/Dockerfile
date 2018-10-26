
FROM php:7

RUN apt-get update && \
    apt-get install -y zlib1g-dev wget

RUN pecl install grpc && \
    docker-php-ext-enable grpc

WORKDIR /usr/bin
RUN php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');" && \
    php -r "if (hash_file('SHA384', 'composer-setup.php') === '93b54496392c062774670ac18b134c3b3a95e5a5e5c8f1a9f115f203b75bf9a129d5daa8ba6a13e2cc8a1da0806388a8') { echo 'Installer verified'; } else { echo 'Installer corrupt'; unlink('composer-setup.php'); } echo PHP_EOL;" && \
    php composer-setup.php && \
    php -r "unlink('composer-setup.php');" && \
    mv composer.phar composer

RUN wget -O phpunit https://phar.phpunit.de/phpunit-7.phar && \
    chmod +x phpunit

