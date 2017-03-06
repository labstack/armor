+++
title = "Template"
description = "Armor template"
[menu.main]
  name = "Template"
  parent = "guide"
+++

Armor template can be used to form a string from variables derived from the HTTP request.
For example, to transfer a url `/users/:name` to `/users/joe` you will use

```sh
"/users/${path:name}"
```

## Available Variables

- `scheme` HTTP scheme (http or https)
- `method` HTTP method
- `uri`	Request URI
- `path` URL path
- `header:<NAME>` Request header
- `path:<NAME>` Path parameter
- `query:<NAME>` Query parameter
- `form:<NAME>` Form parameter

## Supported Plugins

- [`redirect`](/plugin/redirect/#redirect)