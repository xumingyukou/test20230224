package transform

import "net/url"

func Quote(uri string) string {
	return url.QueryEscape(uri)
}
func Unquote(uri string) string {
	res, _ := url.QueryUnescape(uri)
	return res
}
