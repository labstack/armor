+++
title = "Header Plugin"
description = "Header plugin adds / removes response header"
[menu.main]
  name = "Header"
  parent = "plugins"
  weight = 5
+++

Add/remove HTTP response header

## Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `header` | Plugin name
`set` | map | | Set header
`add` | map | | Add header
`del` | array | | Delete header
