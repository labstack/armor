+++
title = "Auto TLS"
description = "Automatic TLS (HTTPS) via Let’s Encrypt"
[menu.main]
  name = "Auto TLS"
  parent = "cookbook"
+++

For demo purpose, I will use domain `test.lab.st` that serves a simple
HTML page.

## Steps

- Have a domain that resolves to an IP address via `A` or `CNAME` DNS record
- [Install Armor]({{< ref "guide.md#installation">}})
- Copy config `config.yaml` to `/etc/armor`
- Copy `index.html` to `/var/www/test`
- Start Armor `armor -c /etc/armor/config.json`
- Browse to `http://test.lab.st`, and in a few seconds TLS certificate will be installed automatically.

## Maintainers

- [vishr](https://github.com/vishr)

## Source

`config.yaml`

```yaml
---
address: :80
tls:
  address: :443
  auto: true
plugins:
- name: https-redirect
hosts:
  test.lab.st:
    plugins:
    - name: static
      root: /var/www/test
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
