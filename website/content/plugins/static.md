+++
title = "Static Plugin"
description = "Static plugin serves static files from a root directory"
[menu.main]
  name = "Static"
  parent = "plugins"
  weight = 4
+++

Serves static files from a provided root directory

## Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `static` | Plugin name
`root` | string | | Root directory from where the static content is served. Required.
`index` | string | `index.html` (default) | Index file for serving a directory
`html5` | bool | `false` (default) | Enable HTML5 mode by forwarding all not-found requests to root so that SPA (single-page application) can handle the routing
`browse` | bool | `false` (default) | Enable directory browsing
