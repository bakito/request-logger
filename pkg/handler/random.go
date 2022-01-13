package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/bakito/request-logger/pkg/common"
	"github.com/gorilla/mux"
)

// RandomCode return a random response code
func RandomCode(w http.ResponseWriter, r *http.Request) {
	code, _ := strconv.Atoi(mux.Vars(r)["code"])
	perc, _ := strconv.ParseFloat(mux.Vars(r)["perc"], 64)

	random := rand.Float64() //nolint:gosec

	delta := random - perc
	if delta <= 0.0 {
		w.WriteHeader(code)
	} else {
		code = 200
	}

	fmt.Printf("%v: %v Code: %v\n", common.HeaderReqNo, w.Header()[common.HeaderReqNo][0], code)
}

// RandomSleep sleep randomly
func RandomSleep(w http.ResponseWriter, r *http.Request) {
	sleep, err := strconv.Atoi(mux.Vars(r)["sleep"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	random := rand.Intn(sleep) //nolint:gosec

	fmt.Printf("%v: %v Sleep: %dms\n", common.HeaderReqNo, w.Header()[common.HeaderReqNo][0], random)
	time.Sleep(time.Duration(random) * time.Millisecond)
	fmt.Printf("%v: %v Sleep: done\n", common.HeaderReqNo, w.Header()[common.HeaderReqNo][0])
}
