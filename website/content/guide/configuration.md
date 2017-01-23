+++
title = "Configuration"
description = "Armor configuration"
[menu.main]
  name = "Configuration"
  parent = "guide"
  weight = 1
+++

Armor accepts configuration in JSON format, command-line option `-c` can be used
to specify a config file, e.g. `armor -c config.json`.

Name | Type | Description
:--- | :--- | :----------
`address` | string | HTTP listen address e.g. `:8080` listens to all IP address on port 8080
`read_timeout` | number | Maximum duration in seconds before timing out read of the request
`write_timeout` | number | Maximum duration before timing out write of the response
`tls` | object | TLS configuration
`plugins` | array | Global plugins
`hosts` | object | Virtual hosts

`tls`

Name | Type | Description
:--- | :--- | :----------
`address` | string | HTTPS listen address. Default value `:80`
`cert_file` | string | Certificate file
`key_file` | string | Key file
`auto` | bool | Enable automatic certificates from https://letsencrypt.org
`cache_dir` | string | Cache directory to store certificates from https://letsencrypt.org. Default value `~/.armor/cache`.

`hosts`

Name | Type | Description
:--- | :--- | :----------
`cert_file` | string | Certificate file
`key_file` | string | Key file
`plugins` | array | Host plugins
`paths` | object | Paths

`paths`

Name | Type | Description
:--- | :--- | :----------
`plugins` | array | Path plugins

## [Plugins]({{< ref "plugins/redirect.md">}})

## Default Configuration

```js
{
	"address": ":8080",
	"plugins": [{
		"name": "logger"
	}, {
		"name": "static",
		"browse": true,
		"root": "."
	}]
}
```

## Sample Configuration

  ```js
{
  "address": ":80",
  "read_timeout": 1200,
  "write_timeout": 1200,
  "tls": {
    "address": ":443",
    "auto": true,
    "cache_dir": "/pool/ingress/cache"
  },
  "plugins": [{
    "name": "logger"
  }, {
    "name": "remove-trailing-slash",
    "redirect_code": 301
  }],
  "hosts": {
    "labstack.com": {
      "paths": {
        "/": {
          "plugins": [{
            "name": "static",
            "root": "/var/www/web",
            "html5": true
          }]
        }
      }
    },
    "api.labstack.com": {
      "paths": {
        "/": {
          "plugins": [{
            "name": "cors"
          }, {
            "name": "proxy",
            "targets": [{
              "url": "http://api.ls"
            }]
          }]
        }
      }
    },
    "armor.labstack.com": {
      "paths": {
        "/": {
          "plugins": [{
            "name": "static",
            "root": "/var/www/armor"
          }]
        }
      }
    },
    "echo.labstack.com": {
      "paths": {
        "/": {
          "plugins": [{
            "name": "static",
            "root": "/var/www/echo"
          }]
        }
      }
    },
    "forum.labstack.com": {
      "paths": {
        "/": {
          "plugins": [{
            "name": "proxy",
            "targets": [{
              "url": "http://forum.ls"
            }]
          }]
        }
      }
    }
  }
}
```