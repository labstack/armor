# [Armor](https://armor.labstack.com) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/labstack/armor/master/LICENSE) [![Build Status](http://img.shields.io/travis/labstack/armor.svg?style=flat-square)](https://travis-ci.org/labstack/echo) [![Join the chat at https://gitter.im/labstack/armor](https://img.shields.io/badge/gitter-join%20chat-brightgreen.svg?style=flat-square)](https://gitter.im/labstack/armor) [![Twitter](https://img.shields.io/badge/twitter-@labstack-55acee.svg?style=flat-square)](https://twitter.com/labstack)

**Uncomplicated HTTP server, supports HTTP/2 and auto TLS**

Armor is written in Go. It is based on the upcoming [Echo](https://github.com/labstack/echo) v3.

## What can it do today?

- Serve HTTP2
- Automatically install TLS certificates from https://letsencrypt.org
- Proxy HTTP and WebSocket requests
- Define virtual hosts with path level routing
- Graceful shutdown
- Limit request body
- Serve static files
- Log requests
- Gzip response
- CORS
- Security
  - XSSProtection
  - ContentTypeNosniff
  - ContentSecurityPolicy
- Add / Remove trailing slash from the URL with option to redirect
- Redirect requests
 - http to https
 - http to https www
 - http to https non www
 - non www to www
 - www to non www
 
Most of the functionality is implemented via `Plugin` interface which makes writing
a custom plugin super easy.

## [Getting Started](https://armor.labstack.com/guide/getting-started)

## What's on the roadmap?

- More command-line options
- More plugins
- More features
- Website
- Code coverage
- Test cases
