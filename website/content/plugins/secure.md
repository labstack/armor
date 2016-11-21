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

Name | Type | Description
:--- | :--- | :----------
`xss_protection` | string | Provides protection against cross-site scripting attack (XSS) by setting the `X-XSS-Protection` header. Default value `1; mode=block`.
`content_type_nosniff` | string | Provides protection against overriding Content-Type header by setting the `X-Content-Type-Options` header. Default value `nosniff`.
`x_frame_options` | string | Can be used to indicate whether or not a browser should be allowed to render a page in a `<frame>`, `<iframe>` or `<object>`. Sites can use this to avoid clickjacking attacks, by ensuring that their content is not embedded into other sites.provides protection against clickjacking. Default value `SAMEORIGIN`. Possible values: `SAMEORIGIN` - The page can only be displayed in a frame on the same origin as the page itself. `DENY` - The page cannot be displayed in a frame, regardless of the site attempting to do so. `ALLOW-FROM uri` - The page can only be displayed in a frame on the specified origin.
`hsts_max_age` | number | Sets the `Strict-Transport-Security` header to indicate how long (in seconds) browsers should remember that this site is only to be accessed using HTTPS. This reduces your exposure to some SSL-stripping man-in-the-middle (MITM) attacks. Default value `0`.
`hsts_exclude_subdomains` | bool | Won't include subdomains tag in the `Strict Transport Security` header, excluding all subdomains from security policy. It has no effect unless HSTSMaxAge is set to a non-zero value. Default value `false`.
`content_security_policy` | string | Sets the `Content-Security-Policy` header providing security against cross-site scripting (XSS), clickjacking and other code injection attacks resulting from execution of malicious content in the  trusted web page context. Default value `""`.
