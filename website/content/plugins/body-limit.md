+++
title = "Body Limit Plugin"
description = "BodyLimit plugin sets the maximum allowed size for a request body"
[menu.main]
  name = "Body Limit"
  parent = "plugins"
  weight = 3
+++

Sets the maximum allowed size for a request body, if the size exceeds the configured
limit, it sends `413 - Request Entity Too Large` response. The body limit is determined
based on both Content-Length request header and actual content read, which makes
it super secure.

## Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `body-limit` | Plugin name
`limit` | string | | Maximum allowed size for a request body, it can be specified as `4x` or `4xB`, where x is one of the multiple from K, M, G, T or P.
