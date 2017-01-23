+++
title = "Trailing Slash Plugins"
description = "Trailing slash plugin adds / removes a trailing slash from the request URI"
[menu.main]
  name = "Trailing Slash"
  parent = "plugins"
  weight = 2
+++

## Add Trailing Slash

Adds a trailing slash to the request URI.

### Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `add-trailing-slash` | Plugin name
`redirect_code` | number | Status code to be used when redirecting the request. Optional, but when provided the request is redirected using this code.

### Remove Trailing Slash

Removes a trailing slash from the request URI.

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `remove-trailing-slash` | Plugin name
`redirect_code` | number | Status code to be used when redirecting the request. Optional, but when provided the request is redirected using this code.
