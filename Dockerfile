FROM alpine:3.7
MAINTAINER Vishal Rana <vr@labstack.com>

ENV VERSION 0.3.7

WORKDIR /

COPY build/armor-${VERSION}_linux-64 /usr/local/bin/armor

ENTRYPOINT ["armor"]