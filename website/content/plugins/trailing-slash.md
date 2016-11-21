+++
title = "Trailing Slash Plugin"
description = "Trailing slash plugin adds / removes a trailing slash from the request URI"
[menu.main]
  name = "TrailingSlash"
  parent = "plugins"
  weight = 2
+++

## Configuration

### `add-trailing-slash`

Adds a trailing slash to the request URI.

### `remove-trailing-slash`

Removes a trailing slash from the request URI.

Name | Type | Description
:--- | :--- | :----------
`redirect_code` | number | Status code to be used when redirecting the request. Optional, but when provided the request is redirected using this code.
