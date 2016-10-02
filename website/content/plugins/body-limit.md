+++
title = "Body Limit"
[menu.side]
  name = "BodyLimit"
  parent = "plugins"
  weight = 3
+++

## BodyLimit Plugin

### `body-limit`

Sets the maximum allowed size for a request body, if the size exceeds the configured
limit, it sends `413 - Request Entity Too Large` response. The body limit is determined
based on both Content-Length request header and actual content read, which makes
it super secure.

### Configuration

Name | Type | Description
:--- | :--- | :----------
`limit` | string | Maximum allowed size for a request body, it can be specified as `4x` or `4xB`, where x is one of the multiple from K, M, G, T or P.
