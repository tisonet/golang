FROM golang:1.10-alpine


ENV GOMAXPROCS=1
ARG LIBRESSL_VERSION=2.6
ARG LIBRDKAFKA_VERSION=0.11.4-RC1B

RUN apk add --no-cache git openssh openssl yajl-dev zlib-dev cyrus-sasl-dev openssl-dev build-base bash gawk sed grep bc coreutils curl tar && \
    apk add libressl${LIBRESSL_VERSION}-libcrypto libressl${LIBRESSL_VERSION}-libssl --update-cache --repository http://nl.alpinelinux.org/alpine/edge/main

RUN curl -L https://github.com/edenhill/librdkafka/archive/v${LIBRDKAFKA_VERSION}.tar.gz | tar xzf - && \
    cd librdkafka-${LIBRDKAFKA_VERSION}/ && \
    ./configure --prefix=/usr && \
    make -j && \
    make install

RUN go get -v github.com/tisonet/golang/openrtb-proxy

ENTRYPOINT ["/go/bin/openrtb-proxy"]