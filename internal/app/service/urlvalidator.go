package service

import (
	"net/url"
	"strings"
)

func ValidateURL(str string) bool {
	path, err := url.ParseRequestURI(str)
	if err == nil && strings.ContainsAny(path.Host, ".:") {
		return true
	}
	return false
}
