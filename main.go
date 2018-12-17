package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/bakito/request-logger/common"
	"github.com/bakito/request-logger/middleware"
	"github.com/gorilla/mux"
)

const (
	defaultPort       int = 8080
	headerTrainReplay     = "Train-Replay"
)

var (
	replayBody        = map[string][]byte{}
	replayContentType = map[string]string{}
)

func main() {

	port := flag.Int("port", defaultPort, "the server port")

	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/echo", echo)
	r.HandleFunc("/echo/{path:.*}", echo)

	r.HandleFunc(`/code/{code:[2,4,5]\d\d}`, responseCode)
	r.HandleFunc(`/code/{code:[2,4,5]\d\d}/{path:.*}`, responseCode)

	r.HandleFunc(`/random/code/{code:[2,4,5]\d\d}/{perc:1|(?:0(?:\.\d*)?)}`, randomCode)
	r.HandleFunc(`/random/code/{code:[2,4,5]\d\d}/{perc:1|(?:0(?:\.\d*)?)}/{path:.*}`, randomCode)

	r.HandleFunc(`/random/sleep/{sleep:\d+}`, randomSleep)
	r.HandleFunc(`/random/sleep/{sleep:\d+}/{path:.*}`, randomSleep)

	r.HandleFunc(`/replay`, replay)
	r.HandleFunc(`/replay/{path:.*}`, replay)

	r.HandleFunc("/{path:.*}", void)

	r.Use(middleware.LogRequest, middleware.CountReqRows)

	log.Printf("Running on port %v ...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), r))
}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, common.Dump(r))
}

func void(w http.ResponseWriter, r *http.Request) {
}

func responseCode(w http.ResponseWriter, r *http.Request) {
	code, _ := strconv.Atoi(mux.Vars(r)["code"])
	w.WriteHeader(code)
	log.Printf("%v: %v Code: %v\n", common.HeaderReqNo, w.Header()[common.HeaderReqNo][0], code)
}

func randomCode(w http.ResponseWriter, r *http.Request) {
	code, _ := strconv.Atoi(mux.Vars(r)["code"])
	perc, _ := strconv.ParseFloat(mux.Vars(r)["perc"], 64)

	random := rand.Float64()

	delta := random - perc
	if delta <= 0.0 {
		w.WriteHeader(code)
	} else {
		code = 200
	}

	log.Printf("%v: %v Code: %v\n", common.HeaderReqNo, w.Header()[common.HeaderReqNo][0], code)
}

func randomSleep(w http.ResponseWriter, r *http.Request) {
	sleep, _ := strconv.Atoi(mux.Vars(r)["sleep"])

	random := rand.Intn(sleep)

	log.Printf("%v: %v Sleep: %dms\n", common.HeaderReqNo, w.Header()[common.HeaderReqNo][0], random)
	time.Sleep(time.Duration(random) * time.Millisecond)
	log.Printf("%v: %v Sleep: done\n", common.HeaderReqNo, w.Header()[common.HeaderReqNo][0])
}

func replay(w http.ResponseWriter, r *http.Request) {
	train := r.Header.Get(headerTrainReplay)

	if train == "true" {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			replayBody[r.RequestURI] = body
		}
		replayContentType[r.RequestURI] = r.Header.Get("Content-Type")
	}

	if t, ok := replayContentType[r.RequestURI]; ok {
		w.Header().Set("Content-Type", t)
	}
	if b, ok := replayBody[r.RequestURI]; ok {
		w.Write(b)
	}
}
