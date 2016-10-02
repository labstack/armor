+++
title = "Proxy"
[menu.side]
  parent = "plugins"
  weight = 3
+++

## Proxy Plugin

### `proxy`

Proxy HTTP and WebSocket requests to upstream servers

### Configuration

Name | Type | Description
:--- | :--- | :----------
`balance` | string | Load balancing technique. Default value `random`. Possible values: `random`, `round-robin`.
`targets` | array | Upstream servers

#### `targets`

Name | Type | Description
:--- | :--- | :----------
`name` | string | Target name
`url` | string | Target url
