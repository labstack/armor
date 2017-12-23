package labstack

import (
	"net"
	"net/http"
	"strings"
)

const (
	rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"
	rfc3339Micro = "2006-01-02T15:04:05.000000Z07:00"
)

// RequestID returns the request ID from the request or response.
func RequestID(r *http.Request, w http.ResponseWriter) string {
	id := r.Header.Get("X-Request-ID")
	if id == "" {
		id = w.Header().Get("X-Request-ID")
	}
	return id
}

// RealIP returns the real IP from the request.
func RealIP(r *http.Request) string {
	ra := r.RemoteAddr
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		ra = strings.Split(ip, ", ")[0]
	} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
		ra = ip
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}
