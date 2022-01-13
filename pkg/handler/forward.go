package handler

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/bakito/request-logger/pkg/common"
)

var skippedHeaders = map[string]bool{
	// skip encoding
	"Accept-Encoding": true,
}

// ForwardFor forward and log request and response
func ForwardFor(target string, disableLogger bool, withTLS bool) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// we need to buffer the body if we want to read it here and send it
		// in the request.
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// you can reassign the body if you need to parse it as multipart
		req.Body = ioutil.NopCloser(bytes.NewReader(body))

		// create a new url from the raw RequestURI sent by the client
		url := joinURL(target, req.URL.Path)
		proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		// We may want to filter some headers, otherwise we could just use a shallow copy
		// proxyReq.Header = req.Header
		proxyReq.Header = make(http.Header)
		for h, val := range req.Header {
			if _, skip := skippedHeaders[h]; !skip {
				proxyReq.Header[h] = val
			}
		}
		if withTLS {
			proxyReq.Header.Set("X-Forwarded-Proto", "https")
		} else {
			proxyReq.Header.Set("X-Forwarded-Proto", "http")
		}

		proxyReq.Header.Set("X-Forwarded-For", readUserIP(req))
		proxyReq.Header.Add("X-Forwarded-Host", req.Host)

		for i, c := range req.Cookies() {
			if i == 0 {
				proxyReq.Header.Set("Cookie", c.String())
			} else {
				proxyReq.Header.Add("Cookie", c.String())
			}
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		}

		httpClient := &http.Client{
			Transport: tr,
		}
		resp, err := httpClient.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		cookies := responseCookies(proxyReq, resp, req)
		for i, c := range cookies {
			if i == 0 {
				resp.Header.Set("Set-Cookie", c.String())
			} else {
				resp.Header.Add("Set-Cookie", c.String())
			}
		}

		if !disableLogger {
			log.Printf("%s (Response): %s\n%s\n", common.HeaderReqNo, w.Header().Get(common.HeaderReqNo), common.DumpResponse(resp))
		}

		for k, vs := range resp.Header {
			for i, v := range vs {
				if i == 0 {
					w.Header().Set(k, v)
				} else {
					w.Header().Add(k, v)
				}
			}
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(respBody)
		w.WriteHeader(resp.StatusCode)
	}
}

func responseCookies(fromReq *http.Request, fromResp *http.Response, to *http.Request) []*http.Cookie {
	host, _, _ := net.SplitHostPort(to.Host)
	cookies := fromResp.Cookies()
	for i := range cookies {
		if fromReq.URL.Host == cookies[i].Domain {
			cookies[i].Domain = host
		}
	}
	return cookies
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}
