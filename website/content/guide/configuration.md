---
title: Configuration
menu:
  side:
    parent: guide
    weight: 2
---

## Configuration

Armor accepts configuration in JSON format, command-line option `-c` can be used
to specify a config file, e.g. `armor -c config.json`.

Name | Type | Description
:--- | :--- | :----------
`address` | string | HTTP listen address e.g. `:8080` listens to all IP address on port 8080
`read_timeout` | number | Maximum duration in seconds before timing out read of the request
`write_timeout` | number | Maximum duration before timing out write of the response
`tls` | object | TLS configuration
`plugins` | object | Global plugins
`hosts` | object | Virtual hosts

### `tls`

Name | Type | Description
:--- | :--- | :----------
`address` | string | HTTPS listen address. Default value `:80`
`cert_file` | string | Certificate file
`key_file` | string | Key file
`auto` | bool | Enable automatic certificates from https://letsencrypt.org
`cache_file` | string | Cache file to store certificates from https://letsencrypt.org. Default value `letsencrypt.cache`.

### `hosts`

Name | Type | Description
:--- | :--- | :----------
`cert_file` | string | Certificate file
`key_file` | string | Key file
`plugins` | object | Host plugins
`paths` | object | Paths

### `paths`

Name | Type | Description
:--- | :--- | :----------
`plugins` | object | Path plugins

### [Plugins]({{< ref "plugins/redirect.md">}})

### Default Configuration

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

### Sample Configuration

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
            "root": "/var/www/web",
            "html5": true
          }
        }
      },
      "blog.labstack.com": {
        "plugins": {
          "static": {
            "root": "/var/www/blog"
          }
        }
      },
      "armor.labstack.com": {
        "plugins": {
          "static": {
            "root": "/var/www/armor"
          }
        }
      },
      "echo.labstack.com": {
        "plugins": {
          "static": {
            "root": "/var/www/echo"
          }
        }
      }
    }
  }
  ```
