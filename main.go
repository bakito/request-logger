package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

const (
	port int = 8080
)

var handler = func() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		log.Printf("Request:\n%s\n\n", (dump))
	})
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(handler())

	log.Printf("Running on port %v ...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", port), r))
}
