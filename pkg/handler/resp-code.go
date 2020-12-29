package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bakito/request-logger/pkg/common"
	"github.com/gorilla/mux"
)

// ResponseCode  return the provided resp code
func ResponseCode(w http.ResponseWriter, r *http.Request) {
	code, err := strconv.Atoi(mux.Vars(r)["code"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	fmt.Printf("%v: %v Code: %v\n", common.HeaderReqNo, w.Header()[common.HeaderReqNo][0], code)
}
