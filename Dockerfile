FROM alpine:3.9

# https://letsencrypt.org
RUN apk add --no-cache ca-certificates

COPY dist/armor_linux_amd64/armor /usr/local/bin

ENTRYPOINT ["armor"]
