+++
title = "Configuration"
description = "Armor configuration"
[menu.main]
  name = "Configuration"
  parent = "guide"
+++

Armor accepts configuration in YAML format, command-line option `-c` can be used
to specify a config file, e.g. `armor -c config.yaml`.

| Name            | Type   | Description                                                             |
| :-------------- | :----- | :---------------------------------------------------------------------- |
| `address`       | string | HTTP listen address e.g. `:8080` listens to all IP address on port 8080 |
| `read_timeout`  | number | Maximum duration in seconds before timing out read of the request       |
| `write_timeout` | number | Maximum duration before timing out write of the response                |
| `tls`           | object | TLS configuration                                                       |
| `plugins`       | array  | Global plugins                                                          |
| `hosts`         | object | Virtual hosts                                                           |

`tls`

| Name            | Type   | Description                                                                                                 |
| :-------------- | :----- | :---------------------------------------------------------------------------------------------------------- |
| `address`       | string | HTTPS listen address. Default value `:80`                                                                   |
| `cert_file`     | string | Certificate file                                                                                            |
| `key_file`      | string | Key file                                                                                                    |
| `auto`          | bool   | Enable automatic certificates from https://letsencrypt.org                                                  |
| `cache_dir`     | string | Cache directory to store certificates from https://letsencrypt.org. Default value `~/.armor/cache`.         |
| `email`         | string | Email optionally specifies a contact email address.                                                         |
| `directory_url` | string | Defines the ACME CA directory endpoint. If empty, LetsEncryptURL is used (acme.LetsEncryptURL).             |
| `secured`       | bool   | If enable, the minimum TLS version is set to 1.2, the ciphers are AEAD and forward secrecy algorithms only. |

`hosts`

| Name        | Type   | Description                                                                                                                 |
| :---------- | :----- | :-------------------------------------------------------------------------------------------------------------------------- |
| `cert_file` | string | Certificate file                                                                                                            |
| `key_file`  | string | Key file                                                                                                                    |
| `plugins`   | array  | Host plugins                                                                                                                |
| `paths`     | object | Paths                                                                                                                       |
| `client_ca` | array  | A list of client CA (certificate authority) certificate encoded as base64 DER. If set client must provide valid certificate |

`paths`

| Name      | Type  | Description  |
| :-------- | :---- | :----------- |
| `plugins` | array | Path plugins |

## [Plugins]({{< ref "plugins/redirect.md">}})

## Default Configuration

```yaml
address: ":8080"
plugins:
- name: logger
- name: static
  browse: true
  root: "."
```

## Sample Configuration

```yaml
address: ":80"
read_timeout: 1200
write_timeout: 1200
tls:
  address: ":443"
  key_file: "/etc/armor/key.pem"
  cert_file: "/etc/armor/cert.pem"
plugins:
- name: logger
- name: cube
  api_key: Y4mYFHJ7jbs1MtVpuGIFirtkvMm9wdJi
- name: https-non-www-redirect
- name: remove-trailing-slash
  redirect_code: 301
hosts:
  labstack.com:
    paths:
      "/":
        plugins:
        - name: proxy
          targets:
          - url: http://web
  api.labstack.com:
    paths:
      "/":
        plugins:
        - name: cors
        - name: proxy
          targets:
          - url: http://api
  armor.labstack.com:
    client_ca_der:
    - "MIIDSzCCAjOgAwI......E/lYx0qGtr0xHQ=="
    paths:
      "/":
        plugins:
        - name: static
          root: "/var/www/armor"
  echo.labstack.com:
    paths:
      "/":
        plugins:
        - name: static
          root: "/var/www/echo"
        - name: redirect
          from: "/recipes*"
          to: "/cookbook${path:*}"
  forum.labstack.com:
    paths:
      "/":
        plugins:
        - name: proxy
          targets:
          - url: http://forum
```
