+++
title = "Static"
[menu.side]
  parent = "plugins"
  weight = 4
+++

## Static Plugin

Serves static files

### `static`

### Configuration

Name | Type | Description
:--- | :--- | :----------
`root` | string | Root directory from where the static content is served. Required.
`index` | string | Index file for serving a directory. Default value `index.html`.
`html5` | bool | Enable HTML5 mode by forwarding all not-found requests to root so that SPA (single-page application) can handle the routing. Default value false.
`browse` | bool | Enable directory browsing. Default value false.
