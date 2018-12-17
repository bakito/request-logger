package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

const (
	defaultPort       int = 8080
	reqNo                 = "Request-No"
	headerReplayTrain     = "REPLAY_TRAIN"
	headerCurrCount       = "CURRENT_COUNT"
)

var (
	counters          sync.Map
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

	r.Use(loggingMiddleware)

	log.Printf("Running on port %v ...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var counter *uint64
		if v, ok := counters.Load(r.RequestURI); ok {
			counter = v.(*uint64)
		} else {
			counter = new(uint64)
			v, _ := counters.LoadOrStore(r.RequestURI, counter)
			counter = v.(*uint64)
		}

		currCount := r.Header.Get(headerCurrCount)

		count := *counter
		if currCount == "" {
			count = atomic.AddUint64(counter, 1)
		}

		log.Printf("%v: %v\n%s\n", reqNo, count, dump(r))

		w.Header().Set(reqNo, fmt.Sprintf("%v", count))
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
	log.Printf("%v: %v Code: %v\n", reqNo, w.Header()[reqNo][0], code)
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

	log.Printf("%v: %v Code: %v\n", reqNo, w.Header()[reqNo][0], code)
}

func randomSleep(w http.ResponseWriter, r *http.Request) {
	sleep, _ := strconv.Atoi(mux.Vars(r)["sleep"])

	random := rand.Intn(sleep)

	log.Printf("%v: %v Sleep: %dms\n", reqNo, w.Header()[reqNo][0], random)
	time.Sleep(time.Duration(random) * time.Millisecond)
	log.Printf("%v: %v Sleep: done\n", reqNo, w.Header()[reqNo][0])
}

func replay(w http.ResponseWriter, r *http.Request) {
	train := r.Header.Get(headerReplayTrain)

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

func dump(r *http.Request) string {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return ""
	}
	return string(dump)
}
