package common

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

const (
	HeaderReqNo = "Request-No"
)

func Dump(r *http.Request) string {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return ""
	}

	req := string(dump)

	if r.RemoteAddr != "" {
		// add remote addre after Host if available, otherwise after the first line
		hi := strings.Index(req, "Host:")
		if hi < 0 {
			hi = 0
		}
		host := req[hi:]
		afterHost := strings.Index(host, "\r\n")
		req = fmt.Sprintf("%s\r\nRemoteAddr: %s%s", req[:hi+afterHost], r.RemoteAddr, req[hi+afterHost:])
	}
	return req
}
