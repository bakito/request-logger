package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/mux"
)

const (
	port int = 8080
)

var counter uint64

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddUint64(&counter, 1)

		log.Printf("Request-No: %v\n%s\n\n", count, dump(r))

		w.Header().Set("Request-No", fmt.Sprintf("%v", count))
		next.ServeHTTP(w, r)
	})
}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, dump(r))
}

func void(w http.ResponseWriter, r *http.Request) {
}

func responseCode(w http.ResponseWriter, r *http.Request) {
	code, _ := strconv.Atoi(mux.Vars(r)["code"])
	w.WriteHeader(code)
}

func main() {
	atomic.AddUint64(&counter, 0)
	r := mux.NewRouter()
	r.HandleFunc("/echo", echo)
	r.HandleFunc("/echo/{path:.*}", echo)
	r.HandleFunc(`/response-code/{code:[2,4,5]\d\d}`, responseCode)
	r.HandleFunc(`/response-code/{code:[2,4,5]\d\d}/{path:.*}`, responseCode)
	r.HandleFunc("/{path:.*}", void)

	r.Use(loggingMiddleware)

	log.Printf("Running on port %v ...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}

func dump(r *http.Request) string {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return ""
	}
	return string(dump)
}
