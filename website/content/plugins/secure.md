+++
title = "Secure Plugin"
description = "Secure plugin provides protection against web attacks"
[menu.main]
  name = "Secure"
  parent = "plugins"
  weight = 3
+++

Secure plugin provides protection against cross-site scripting (XSS) attack,
content type sniffing, clickjacking, insecure connection and other code injection
attacks.

## Configuration

Name | Type | Value | Description
:--- | :--- | :--- | :----------
`name` | string | `secure` | Plugin name
`xss_protection` | string | `1; mode=block` (default) | Provides protection against cross-site scripting attack (XSS) by setting the `X-XSS-Protection` header
`content_type_nosniff` | string | `nosniff` (default) | Provides protection against overriding Content-Type header by setting the `X-Content-Type-Options` header
`x_frame_options` | string | `SAMEORIGIN` (default) | Can be used to indicate whether or not a browser should be allowed to render a page in a `<frame>`, `<iframe>` or `<object>`. Sites can use this to avoid clickjacking attacks, by ensuring that their content is not embedded into other sites.provides protection against clickjacking. Possible values: `SAMEORIGIN` - The page can only be displayed in a frame on the same origin as the page itself. `DENY` - The page cannot be displayed in a frame, regardless of the site attempting to do so. `ALLOW-FROM uri` - The page can only be displayed in a frame on the specified origin.
`hsts_max_age` | number | `0` (defaul) | Sets the `Strict-Transport-Security` header to indicate how long (in seconds) browsers should remember that this site is only to be accessed using HTTPS. This reduces your exposure to some SSL-stripping man-in-the-middle (MITM) attacks.
`hsts_exclude_subdomains` | bool | `false` (default) | Won't include subdomains tag in the `Strict Transport Security` header, excluding all subdomains from security policy. It has no effect unless HSTSMaxAge is set to a non-zero value.
`content_security_policy` | string | `""` (default) | Sets the `Content-Security-Policy` header providing security against cross-site scripting (XSS), clickjacking and other code injection attacks resulting from execution of malicious content in the  trusted web page context.
