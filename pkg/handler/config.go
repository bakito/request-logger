package handler

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/bakito/request-logger/pkg/conf"
)

// ConfigReplay replay from config
func ConfigReplay(resp conf.Response) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if resp.ContentType != "" {
			w.Header().Set("Content-Type", resp.ContentType)
		} else {
			if r.Header.Get("Accept") != "" {
				w.Header().Set("Content-Type", r.Header.Get("Accept"))
			} else {
				w.Header().Set("Content-Type", "text/plain")
			}
		}

		var data []byte
		var err error
		if resp.BodyFile != "" {
			data, err = ioutil.ReadFile(resp.BodyFile)
		} else {
			data = []byte(resp.Body)
		}
		if err == nil {
			_, err = w.Write(data)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// ConfigLogBody log body from config file
func ConfigLogBody(lb conf.LogBody) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(headerLogBodyLength, strconv.FormatBool(lb.LineLength))
		if lb.ResponseCode != 0 {
			w.WriteHeader(lb.ResponseCode)
		}
		LogBody(w, r)
	}
}
