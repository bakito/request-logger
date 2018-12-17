package middleware

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/bakito/request-logger/common"
)

const (
	headerCurrCount    = "Current-Count"
	headerCountReqRows = "Count-Request-Rows"
	headerTotalReqRows = "Total-Request-Rows"
)

var (
	reqCounters    sync.Map
	reqRowCounters sync.Map
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var counter *uint64
		if v, ok := reqCounters.Load(r.RequestURI); ok {
			counter = v.(*uint64)
		} else {
			v, _ := reqCounters.LoadOrStore(r.RequestURI, new(uint64))
			counter = v.(*uint64)
		}

		currCount := r.Header.Get(headerCurrCount)

		count := *counter
		if currCount == "" {
			count = atomic.AddUint64(counter, 1)
		}

		log.Printf("%v: %v\n%s\n", common.HeaderReqNo, count, common.Dump(r))

		w.Header().Set(common.HeaderReqNo, fmt.Sprintf("%v", count))
		next.ServeHTTP(w, r)
	})
}

func CountReqRows(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cnt := r.Header.Get(headerCountReqRows)

		if cnt == "true" {
			reqRowCounters.Store(r.RequestURI, new(uint64))
		} else if v, ok := reqRowCounters.Load(r.RequestURI); ok {

			var lines uint64
			scanner := bufio.NewScanner(r.Body)
			for scanner.Scan() {
				lines++
			}

			counter := v.(*uint64)

			allLines := atomic.AddUint64(counter, uint64(lines))

			w.Header().Set(headerTotalReqRows, fmt.Sprintf("%v", allLines))
		}

		next.ServeHTTP(w, r)
	})
}
