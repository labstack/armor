FROM alpine:3.7

# https://letsencrypt.org
RUN apk add --no-cache ca-certificates

COPY dist/linux_amd64/armor /usr/local/bin

ENTRYPOINT ["armor"]
