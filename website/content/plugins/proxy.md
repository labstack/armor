+++
title = "Proxy Plugin"
description = "Proxy plugin proxies HTTP and WebSocket requests to upstream servers"
[menu.side]
  name = "Proxy"
  parent = "plugins"
  weight = 3
+++

## `proxy`

Proxy HTTP and WebSocket requests to upstream servers

## Configuration

Name | Type | Description
:--- | :--- | :----------
`balance` | string | Load balancing technique. Default value `random`. Possible values: `random`, `round-robin`.
`targets` | array | Upstream servers

### `targets`

Name | Type | Description
:--- | :--- | :----------
`name` | string | Target name
`url` | string | Target url
