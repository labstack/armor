+++
title = "CORS Plugin"
description = "CORS plugin gives web servers cross-domain access controls, which enable secure cross-domain data transfers."
[menu.main]
  name = "CORS"
  parent = "plugins"
  weight = 4
+++

[CORS](http://www.w3.org/TR/cors) gives web servers cross-domain access controls,
which enable secure cross-domain data transfers.

## Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | array | `cors` | Plugin name
`allow_origins` | array | `["*"]` (default) | Defines a list of origins that may access the resource
`allow_methods` | array | `["HEAD", "PUT", "PATCH", "POST", "DELETE"]` (default) | Defines a list methods allowed when accessing the resource. This is used in response to a preflight request.
`allow_headers` | array | `[]` (default) | Defines a list of request headers that can be used when making the actual request. This in response to a preflight request.
`allow_credentials` | bool | `false` (default) | Indicates whether or not the response to the request can be exposed when the credentials flag is true. When used as part of a response to a preflight request, this indicates whether or not the  actual request can be made using credentials.
`expose_headers` | array | `[]` (default) | Defines a whitelist headers that clients are allowed to access
`max_age` | number | `0` (default) | Indicates how long (in seconds) the results of a preflight request can be cached
