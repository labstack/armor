# Build image
FROM golang:alpine AS builder
MAINTAINER Vishal Rana <vr@labstack.com>

COPY . /go/src/github.com/labstack/armor

WORKDIR /go/src/github.com/labstack/armor

RUN set -x \
    && export CGO_ENABLED=0 \
    && go build -v -o /go/bin/armor cmd/armor/main.go

# Executable image
FROM scratch

WORKDIR /

COPY --from=builder /go/bin/armor /armor
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/armor"]