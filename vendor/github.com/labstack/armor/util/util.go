package util

import "strings"

// StripPort strips port from the host.
// https://golang.org/pkg/net/url/#URL.Hostname
func StripPort(host string) string {
	colon := strings.IndexByte(host, ':')
	if colon == -1 {
		return host
	}
	if i := strings.IndexByte(host, ']'); i != -1 {
		return strings.TrimPrefix(host[:i], "[")
	}
	return host[:colon]
}
