package middleware

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/bakito/request-logger/pkg/common"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const (
	headerTotalReqRows   = "Total-Request-Rows"
	headerCurrentReqRows = "Current-Request-Rows"
)

var (
	reqRowCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_logger_request_body_row_count",
		Help: "The total count rows in request body for a path",
	}, []string{"path"})
	currCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_logger_request_count",
		Help: "The count of requests by path",
	}, []string{"path"})
)

func init() {
	prometheus.MustRegister(currCount, reqRowCount)
}

// CountRequests count the number of requests per path
func CountRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := inc(currCount, r)
		w.Header().Set(common.HeaderReqNo, fmt.Sprintf("%v", count))
		next.ServeHTTP(w, r)
	})
}

// LogRequest logging middleware
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s\n%s\n", common.HeaderReqNo, w.Header().Get(common.HeaderReqNo), common.DumpRequest(r))
		next.ServeHTTP(w, r)
	})
}

// CountReqRows row counting middleware
func CountReqRows(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var lines float64

		bodyBytes := common.GetBody(r)

		scanner := bufio.NewScanner(bytes.NewReader(bodyBytes))
		for scanner.Scan() {
			lines++
		}

		allLines := add(reqRowCount, r, lines)
		w.Header().Set(headerTotalReqRows, fmt.Sprintf("%v", allLines))
		w.Header().Set(headerCurrentReqRows, fmt.Sprintf("%v", lines))

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
	_ = c.Write(pb)
	return pb.GetCounter().GetValue()
}
