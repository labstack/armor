---
title: Index
---

# Armor

**Uncomplicated HTTP server, supports HTTP/2 and auto TLS**

## What can it do today?

- Serve HTTP/2
- Automatically install TLS certificates from https://letsencrypt.org
- Proxy HTTP and WebSocket requests
- Define virtual hosts with path level routing
- Graceful shutdown
- Limit request body
- Serve static files
- Log requests
- Gzip response
- Cross-origin Resource Sharing (CORS)
- Security
  - XSSProtection
  - ContentTypeNosniff
  - ContentSecurityPolicy
  - HTTP Strict Transport Security (HSTS)
- Add / Remove trailing slash from the URL with option to redirect
- Redirect requests
 - http to https
 - http to https www
 - http to https non www
 - non www to www
 - www to non www

Most of the functionality is implemented via `Plugin` interface which makes writing
a custom plugin super easy.

### [Getting Started]({{< ref "guide/getting-started.md">}})
