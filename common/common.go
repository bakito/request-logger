package common

import (
	"net/http"
	"net/http/httputil"
)

const (
	HeaderReqNo = "Request-No"
)

func Dump(r *http.Request) string {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return ""
	}
	return string(dump)
}
