package helix

import (
	"strings"
)

func ETagMatch(header string, etag string) bool {
	if len(etag) == 0 {
		return true
	}
	if len(header) == 0 {
		return true
	}
	val := strings.Split(header, ",")
	for _, v := range val {
		v = strings.TrimSpace(v)
		if v == "*" || v == etag {
			return true
		}
	}
	return false
}

func ETagNoneMatch(header string, etag string) bool {
	if len(etag) == 0 {
		return true
	}
	if len(header) == 0 {
		return true
	}
	val := strings.Split(header, ",")
	for _, v := range val {
		v = strings.TrimSpace(v)
		if v != "*" && v != etag {
			return true
		}
	}
	return false
}
