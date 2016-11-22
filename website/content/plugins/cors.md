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

Name | Type | Description
:--- | :--- | :----------
`allow_origins` | array | Defines a list of origins that may access the resource. Default value ["*"].
`allow_methods` | array | Defines a list methods allowed when accessing the resource. This is used in response to a preflight request. Default value ["HEAD", "PUT", "PATCH", "POST", "DELETE"].
`allow_headers` | array | Defines a list of request headers that can be used when making the actual request. This in response to a preflight request. Default value [].
`allow_credentials` | bool | Indicates whether or not the response to the request can be exposed when the credentials flag is true. When used as part of a response to a preflight request, this indicates whether or not the  actual request can be made using credentials. Default value `false`.
`expose_headers` | array | Defines a whitelist headers that clients are allowed to access. Default value [].
`max_age` | number | Indicates how long (in seconds) the results of a preflight request can be cached. Default value `0`.
