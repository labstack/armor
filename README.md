<a href="https://armor.labstack.com"><img height="80" src="https://labstack.com/images/armor-logo.svg"></a>

[![Build Status](http://img.shields.io/travis/labstack/armor.svg?style=flat-square)](https://travis-ci.org/labstack/armor)
[![Forum](https://img.shields.io/badge/community-forum-00afd1.svg?style=flat-square)](https://forum.labstack.com)
[![Twitter](https://img.shields.io/badge/twitter-@labstack-55acee.svg?style=flat-square)](https://twitter.com/labstack)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/labstack/armor/master/LICENSE)

## What can it do today?

- Serve HTTP/2
- Automatically install TLS certificates from https://letsencrypt.org
- Proxy HTTP and WebSocket requests
- Define virtual hosts with path level routing
- Graceful shutdown
- Limit request body
- Serve static files
- Log requests
- Gzip response
- Cross-origin Resource Sharing (CORS)
- Security
  - XSSProtection
  - ContentTypeNosniff
  - ContentSecurityPolicy
  - HTTP Strict Transport Security (HSTS)
- Add / Remove trailing slash from the URL with option to redirect
- Redirect requests
 - http to https
 - http to https www
 - http to https non www
 - non www to www
 - www to non www

Most of the functionality is implemented via `Plugin` interface which makes writing
a custom plugin super easy.

## [Get Started](https://armor.labstack.com/guide)

## What's on the roadmap?

- [x] Website
- [ ] Code coverage
- [ ] Test cases
