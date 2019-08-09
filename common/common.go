package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

const (
	// HeaderReqNo http header
	HeaderReqNo = "Request-No"
)

// Dump dump the request into a string
func Dump(r *http.Request) string {
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

// GetBody get the content of the body, closes the read one and adds a new read closer to the request
func GetBody(r *http.Request) []byte {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close() //  must close
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}
