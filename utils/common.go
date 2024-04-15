package utils

import "net/url"

func Or(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	for _, str := range s[:len(s)-1] {
		if str != "" {
			return str
		}
	}
	return s[len(s)-1]
}

func DecodeURIComponent(str string) string {
	result, _ := url.QueryUnescape(str)

	return result
}
