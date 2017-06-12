+++
title = "Redirect Plugins"
description = "Redirect plugin redirects HTTP requests"
[menu.main]
  name = "Redirect"
  parent = "plugins"
  weight = 2
+++

## Redirect

Redirects http requests base on *from* and *to* URL.

### Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `redirect` | Plugin name
`from` | string | | Redirect from URI
`to` | string (template) | | Redirect to URI
`code` | number | `301` (default) | Redirect code

*Example*

```yaml
name: redirect
from: "/recipes*"
to: "/cookbook${path:*}"
```

## HTTPS Redirect

Redirects http requests to https. For example, http://labstack.com will be redirected
to https://labstack.com.

### Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `https-redirect` | Plugin name
`code` | number | `301` (default) | Redirect code

## HTTPS WWW Redirect

Redirects http requests to www https. For example, http://labstack.com will be redirected to https://www.labstack.com.

### Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `https-www-redirect` | Plugin name
`code` | number | `301` (default) | Redirect code

## HTTPS Non-WWW Redirect

Redirects http requests to https non-www. For example, http://www.labstack.com will
be redirect to https://labstack.com.

### Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `https-non-www-redirect` | Plugin name
`code` | number | `301` (default) | Redirect code

## Non-WWW Redirect

Redirects www requests to non-www. For example, http://www.labstack.com will be
redirected to http://labstack.com.

### Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `non-www-redirect` | Plugin name
`code` | number | `301` (default) | Redirect code

## WWW Redirect

Redirects non-www requests to www.
For example, http://labstack.com will be redirected to http://www.labstack.com.

### Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `www-redirect` | Plugin name
`code` | number | `301` (default) | Redirect code
