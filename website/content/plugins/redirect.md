+++
title = "Redirect Plugin"
description = "Redirect plugin redirects HTTP requests"
[menu.main]
  name = "Redirect"
  parent = "plugins"
  weight = 2
+++

## `https-redirect`

Redirects http requests to https. For example, http://labstack.com will be redirected
to https://labstack.com.

## `https-www-redirect`

Redirects http requests to www https. For example, http://labstack.com will be redirected to https://www.labstack.com.

## `https-non-www-redirect`

Redirects http requests to https non www. For example, http://www.labstack.com will
be redirect to https://labstack.com.

## `non-www-redirect`

Redirects www requests to non www. For example, http://www.labstack.com will be
redirected to http://labstack.com.

## `www-redirect`

Redirects non www requests to www.
For example, http://labstack.com will be redirected to http://www.labstack.com.

## Configuration

Name | Type | Description
:--- | :--- | :----------
`code` | number | Redirect code. Default value `301` (Moved Permanently).
