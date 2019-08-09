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

	"github.com/bakito/request-logger/conf"

	"github.com/bakito/request-logger/common"
	"github.com/bakito/request-logger/middleware"
	"github.com/gorilla/mux"
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
	disableLogger := flag.Bool("disableLogger", false, "Disable the logger middleware")
	configFile := flag.String("config", "", "The path of a config file")

	flag.Parse()

	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	var config *conf.Conf
	var err error

	if configFile != nil && *configFile != "" {
		config, err = conf.GetConf(*configFile)
		if err != nil {
			log.Fatalf("Error reading config %s: %v", *configFile, err)
			return
		}
	}

	functions := make(map[string]func(w http.ResponseWriter, r *http.Request))
	var paths []string

	if config != nil {

		for _, path := range config.Echo {
			functions[path] = echo
			paths = append(paths, path)
		}

		for _, path := range config.LogBody {
			functions[path] = logBody
			paths = append(paths, path)
		}

		for _, resp := range config.Replay {
			functions[resp.Path] = configReplay(resp)
			paths = append(paths, resp.Path)
		}

		common.SortPaths(paths)

		log.Printf("Serving custom config from '%s'", *configFile)
		for _, p := range paths {
			r.HandleFunc(p, functions[p])
		}

	} else {
		r.HandleFunc("/echo", echo)
		r.HandleFunc("/echo/{path:.*}", echo)

		r.HandleFunc("/body", logBody)
		r.HandleFunc("/body/{path:.*}", logBody)

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
	if !*disableLogger {
		r.Use(middleware.LogRequest)
	}
	if *countRequestRows {
		r.Use(middleware.CountReqRows)
	}

	log.Printf("Running on port %v ...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), r))
}
func configReplay(resp conf.Response) func(w http.ResponseWriter, r *http.Request) {
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

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprint(w, common.Dump(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logBody(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err == nil {
		log.Print(string(body))
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func void(w http.ResponseWriter, r *http.Request) {
}

func responseCode(w http.ResponseWriter, r *http.Request) {
	code, err := strconv.Atoi(mux.Vars(r)["code"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	sleep, err := strconv.Atoi(mux.Vars(r)["sleep"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
		defer r.Body.Close()
	}

	if t, ok := replayContentType[r.RequestURI]; ok {
		w.Header().Set("Content-Type", t)
	}
	if b, ok := replayBody[r.RequestURI]; ok {
		_, err := w.Write(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
