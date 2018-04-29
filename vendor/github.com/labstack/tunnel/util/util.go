package util

import (
	"strconv"
	"strings"
)

func SplitHostPort(address string) (host string, port int, err error) {
	parts := strings.Split(address, ":")
	if len(parts) == 1 {
		port, err = strconv.Atoi(parts[0])
		if err != nil {
			return
		}
	} else if len(parts) == 2 {
		host = parts[0]
		port, err = strconv.Atoi(parts[1])
		if err != nil {
			return
		}
	}
	return
}
