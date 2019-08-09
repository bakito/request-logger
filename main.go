package main

import (
	"fmt"
	"github.com/bakito/request-logger/conf"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/bakito/request-logger/common"
	"github.com/bakito/request-logger/middleware"
	"github.com/gorilla/mux"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	countRequestRows := flag.Bool("countRequestRows", true, "Enable or disable the request row count")

	flag.Parse()

	r := mux.NewRouter()

	config := conf.GetConf()
	if config != nil {

		for _, path := range config.Echo {
			r.HandleFunc(path, echo)
		}
		for _, path := range config.EchoBody {
			r.HandleFunc(path, body)
		}

		for _, resp := range config.Replay {
			r.HandleFunc(resp.Path, func(w http.ResponseWriter, r *http.Request) {
				if resp.ContentType != "" {
					w.Header().Set("Content-Type", resp.ContentType)
				} else {
					w.Header().Set("Content-Type", "text/plain")
				}
				_, err := w.Write([]byte(resp.Content))
				if err != nil {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			})
		}
	} else {
		r.Handle("/metrics", promhttp.Handler())

		r.HandleFunc("/echo", echo)
		r.HandleFunc("/echo/{path:.*}", echo)

		r.HandleFunc("/body", body)
		r.HandleFunc("/body/{path:.*}", body)

		r.HandleFunc(`/code/{code:[2,4,5]\d\d}`, responseCode)
		r.HandleFunc(`/code/{code:[2,4,5]\d\d}/{path:.*}`, responseCode)

		r.HandleFunc(`/random/code/{code:[2,4,5]\d\d}/{perc:1|(?:0(?:\.\d*)?)}`, randomCode)
		r.HandleFunc(`/random/code/{code:[2,4,5]\d\d}/{perc:1|(?:0(?:\.\d*)?)}/{path:.*}`, randomCode)

		r.HandleFunc(`/random/sleep/{sleep:\d+}`, randomSleep)
		r.HandleFunc(`/random/sleep/{sleep:\d+}/{path:.*}`, randomSleep)

		r.HandleFunc(`/replay`, replay)
		r.HandleFunc(`/replay/{path:.*}`, replay)

		r.HandleFunc("/{path:.*}", void)
	}
	r.Use(middleware.LogRequest)
	if *countRequestRows {
		r.Use(middleware.CountReqRows)
	}

	log.Printf("Running on port %v ...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), r))
}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, common.Dump(r))
}

func body(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		_, _ = w.Write(body)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
		_, _ = w.Write(b)
	}
}
