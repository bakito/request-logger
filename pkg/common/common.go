package common

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

const (
	// HeaderReqNo http header
	HeaderReqNo = "Request-No"
)

// DumpRequest dump the request into a string
func DumpRequest(r *http.Request) string {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return ""
	}

	req := string(dump)

	if r.RemoteAddr != "" {
		// add remote address after Host if available, otherwise after the first line
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

// DumpResponse dump the response into a string
func DumpResponse(resp *http.Response) string {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return ""
	}

	req := string(dump)

	return req
}

// GetBody get the content of the body, closes the read one and adds a new read closer to the request
func GetBody(r *http.Request) []byte {
	bodyBytes, _ := io.ReadAll(r.Body)
	_ = r.Body.Close() //  must close
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}
