package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"sync/atomic"

	"github.com/gorilla/mux"
)

const (
	port int = 8080
)

var counter uint64

var handler = func() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&counter, 1)
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Printf("Request %v:\n%s\n\n", counter, dump)
	})
}

func main() {
	atomic.AddUint64(&counter, 0)
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(handler())
	log.Printf("Running on port %v ...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
