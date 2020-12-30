package handler

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
		url := fmt.Sprintf("%s/%s", target, req.URL.Path)
		proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))

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

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient := &http.Client{Transport: tr}
		resp, err := httpClient.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if !disableLogger {
			log.Printf("%s (Response): %s\n%s\n", common.HeaderReqNo, w.Header().Get(common.HeaderReqNo), common.DumpResponse(resp))
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)
		for k, vs := range resp.Header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		w.Write(respBody)
	}
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
