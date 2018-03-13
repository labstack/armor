FROM alpine:3.7

# https://letsencrypt.org
RUN apk add --no-cache ca-certificates

COPY build/armor-*-linux-amd64 /usr/local/bin/armor

ENTRYPOINT ["armor"]