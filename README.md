# Armor - Simple HTTP server, supports HTTP/2 and auto TLS

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/labstack/armor/master/LICENSE) [![Build Status](http://img.shields.io/travis/labstack/armor.svg?style=flat-square)](https://travis-ci.org/labstack/echo) [![Join the chat at https://gitter.im/labstack/armor](https://img.shields.io/badge/gitter-join%20chat-brightgreen.svg?style=flat-square)](https://gitter.im/labstack/armor) [![Twitter](https://img.shields.io/badge/twitter-@labstack-55acee.svg?style=flat-square)](https://twitter.com/labstack)

Armor is a simple HTTP server written in Go. It is based on the upcoming [Echo](https://github.com/labstack/echo) v3.

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
 - HTTP to HTTPS
 - HTTP to HTTPS WWW
 - HTTP to HTTPS non WWW
 - Non WWW to WWW
 - WWW to non WWW

Most of the functionality is implemented via `Plugin` interface which makes writing
a custom plugin super easy.

## Getting Started

### Installation

- Download the latest armor release for your platform from https://github.com/labstack/armor/releases
- Copy the armor binary to somewhere on the `PATH` so that it can be executed e.g. `/usr/local/bin` or `%PATH%`

### Executing

Open a terminal and type `armor`

```sh
❯ armor

 _______  ______    __   __  _______  ______
|   _   ||    _ |  |  |_|  ||       ||    _ |
|  |_|  ||   | ||  |       ||   _   ||   | ||
|       ||   |_||_ |       ||  | |  ||   |_||_
|       ||    __  ||       ||  |_|  ||    __  |
|   _   ||   |  | || ||_|| ||       ||   |  | |
|__| |__||___|  |_||_|   |_||_______||___|  |_|

                                      v0.1.1

Simple HTTP server, supports HTTP/2 and auto TLS
      https://github.com/labstack/armor
___________________O/___________________________
                   O\

 ⇛ http server started on :8080
```

This starts armor on address `:8080` serving the current directory listing using a
default config. Browse to http://localhost:8080 to see the listing.

Armor can also be run using Docker `docker run labstack/armor`

### Configuration

Armor accepts configuration in JSON format, command-line option `-c` can be used
to specify a config file, e.g. `armor -c config.json`.

#### Default Config

```js
{
  "address": ":8080",
  "plugins": {
    "logger": {},
    "static": {
      "browse": true,
      "root": "."
    }
  }
}
```

#### General

- `address`(string): HTTP listen address e.g. `:8080` listens to all IP address on port 8080
- `tls`(object): TLS configuration
  - `address`(string): HTTPS listen address
  - `cert_file`(string): Certificate file
  - `key_file`(string): Key file
  - `auto`(bool): Enable automatic certificates from https://letsencrypt.org
	- `cache_file`(string): Cache file to store certificates from https://letsencrypt.org. Optional. Default value letsencrypt.cache.
- `read_timeout`(number - in seconds): Maximum duration before timing out read of the request
- `write_timeout`(number - in seconds): Maximum duration before timing out write of the response
- `plugins`(object): Global level plugins
- `hosts`(object): Virtual hosts
  - `cert_file`(string): Certificate file
  - `key_file`(string): Key file
  - `plugins`(object): Host level plugins
	- `paths`(object): Paths
    - `plugins`(object) Path level plugins

#### Plugins

- `body_limit`(object): https://echo.labstack.com/middleware/body-limit
  - `limit`(string)
- `cors`(object): https://echo.labstack.com/middleware/cors
  - `allow_origins`([]string)
  - `allow_methods`([]string)
  - `allow_headers`([]string)
  - `allow_credentials`(bool)
  - `expose_headers`([]string)
  - `max_age`(number)
- `gzip`(object): https://echo.labstack.com/middleware/gzip
  - `level`(number)
- `header`(object): Add / remove response header.
  - `set`(string): Set header
  - `add`(string): Add header
  - `del`(string): Delete header
- `logger`(object): https://echo.labstack.com/middleware/logger
  - `format`(string)
- `proxy`(object)
  - `balance`(string)
- `https-redirect`(object): https://echo.labstack.com/middleware/redirect
  - `code`(number)
- `https-www-redirect`(object)
  - `code`(number)
- `https-non-www-redirect`(object)
  - `code`(number)
- `non-www-redirect`(object)
  - `code`(number)
- `www-redirect`(object)
  - `code`(number)
- `add-trailing-slash`(object): https://echo.labstack.com/middleware/trailing-slash
  - `redirect_code`(number)
- `remove-trailing-slash`(object)
  - `redirect_code`(number)
- `secure`(object): https://echo.labstack.com/middleware/secure
  - `xss_protection`(string)
  - `content_type_nosniff`(string)
  - `x_frame_options`(string)
  - `hsts_max_age`(number)
  - `hsts_exclude_subdomains`(bool)
  - `content_security_policy`(string)
- `static`(object): https://echo.labstack.com/middleware/static
  - `root`(string)
  - `index`(string)
  - `html5`(bool)
  - `browse`(bool)

  #### Sample Configuration

  ```js
  {
    "address": ":80",
    "tls": {
      "auto": true
    },
    "plugins": {
      "https-redirect": {},
      "remove-trailing-slash": {
        "redirect_code": 301
      },
      "logger": {},
      "gzip": {}
    },
    "hosts": {
      "api.labstack.com": {
        "plugins": {
            "cors": {},
            "proxy": {
                "targets": [{
                    "url": "http://api.ls"
                }]
            }
        }
      },
      "labstack.com": {
        "plugins": {
          "non-www-redirect": {},
          "static": {
            "root": "/srv/web",
            "html5": true
          }
        }
      },
      "blog.labstack.com": {
        "plugins": {
          "static": {
            "root": "/srv/blog"
          }
        }
      }
    }
  }
  ```

## What's on the roadmap?

- More command-line options
- More plugins
- More features
- Website
- Code coverage
- Test cases
