package util

import (
	"strings"

	"github.com/labstack/gommon/random"
)

func ID() string {
	return random.String(8, random.Lowercase, random.Numeric)
}

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
