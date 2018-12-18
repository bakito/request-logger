package middleware

import (
	"bufio"
	"fmt"
	"log"
	"net/http"

	"github.com/namsral/flag"

	"github.com/bakito/request-logger/common"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const (
	headerCurrCount    = "Current-Count"
	headerTotalReqRows = "Total-Request-Rows"
)

var (
	reqRowCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_body_row_count",
		Help: "The current count rows in request body",
	}, []string{"path"})
	currCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_current_count",
		Help: "The current count of requests by path",
	}, []string{"path"})
	dumpRequest *bool
)

func init() {
	prometheus.MustRegister(currCount, reqRowCount)
	dumpRequest = flag.Bool("dumpRequest", true, "Enable or disable the request dump")
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		count := inc(currCount, r)

		if *dumpRequest {
			log.Printf("%v: %v\n%s\n", common.HeaderReqNo, count, common.Dump(r))
		} else {
			reqURI := r.RequestURI
			if reqURI == "" {
				reqURI = r.URL.RequestURI()
			}
			log.Printf("%v: %v\n%s %s HTTP/%d.%d\n", common.HeaderReqNo, count, valueOrDefault(r.Method, "GET"),
				reqURI, r.ProtoMajor, r.ProtoMinor)
		}

		w.Header().Set(common.HeaderReqNo, fmt.Sprintf("%v", count))
		next.ServeHTTP(w, r)
	})
}

func CountReqRows(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var lines float64
		scanner := bufio.NewScanner(r.Body)
		for scanner.Scan() {
			lines++
		}

		allLines := add(reqRowCount, r, lines)
		w.Header().Set(headerTotalReqRows, fmt.Sprintf("%v", allLines))

		next.ServeHTTP(w, r)
	})
}

func inc(cVec *prometheus.CounterVec, r *http.Request) float64 {
	return add(cVec, r, 1)
}

func add(cVec *prometheus.CounterVec, r *http.Request, value float64) float64 {
	c := cVec.WithLabelValues(r.RequestURI)
	c.Add(value)
	pb := &dto.Metric{}
	c.Write(pb)
	return pb.GetCounter().GetValue()
}

func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}
