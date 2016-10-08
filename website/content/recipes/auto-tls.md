+++
title = "Auto TLS"
description = "Automatic TLS (HTTPS) via Let’s Encrypt"
[menu.side]
  name = "Auto TLS"
  parent = "recipes"
  weight = 1
+++

## Automatic TLS (HTTPS) via [Let’s Encrypt](https://letsencrypt.org)

For demo purpose, I will use domain `test.lab.st` that serves a simple
HTML page.

### Steps

- Have a domain that resolves to an IP address via `A` or `CNAME` DNS record
- [Install Armor]({{< ref "guide/getting-started.md#installation">}})
- Copy config `config.json` to `/etc/armor`
- Copy `index.html` to `/var/www/test`
- Start Armor `armor -c /etc/armor/config.json`
- Browse to `http://test.lab.st`, and in a few seconds TLS certificate will
be installed automatically.

### Source

`config.json`

```js
{
  "tls": {
    "auto": true,
    "cache_file": "/var/www/le.cache"
  },
  "plugins": {
    "https-redirect": {}
  },
  "hosts": {
    "test.lab.st": {
      "plugins": {
        "static": {
          "root": "/var/www/test"
        }
      }
    }
  }
}
```

`index.html`

```html
<!doctype html>
<html lang="en-us">
<head>
  <title>Armor</title>
</head>
<body>
  <h1>Welcome to Armor!</h1>
</body>
</html>
```
