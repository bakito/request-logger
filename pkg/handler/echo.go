package handler

import (
	"fmt"
	"net/http"

	"github.com/bakito/request-logger/pkg/common"
)

// Echo the request
func Echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprint(w, common.DumpRequest(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
