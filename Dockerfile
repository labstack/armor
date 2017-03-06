FROM alpine:edge
MAINTAINER Vishal Rana <vr@labstack.com>

ENV VERSION 0.2.8

# https://letsencrypt.org
RUN apk add --no-cache ca-certificates

# TODO: version variable
COPY build/armor-${VERSION}_linux-64 /usr/local/bin/armor

ENTRYPOINT ["armor"]
