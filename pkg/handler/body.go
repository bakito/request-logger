package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	headerLogBodyLength = "Log-Body-Length"
	headerBodyAsString  = "Log-Body-As-String"
)

// LogBody log the body
func LogBody(logAsString bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		length := r.Header.Get(headerLogBodyLength) == "true"
		logAsString = logAsString || r.Header.Get(headerBodyAsString) == "true"

		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		body, err := ioutil.ReadAll(r.Body)
		defer func() { _ = r.Body.Close() }()
		if err == nil {
			r := bufio.NewReader(bytes.NewReader(body))
			lineNbr := 0
			for {
				line, _, err := r.ReadLine()

				if err == io.EOF {
					break
				} else if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					break
				}

				if length {
					if lineNbr == 0 {
						fmt.Printf("%10d | %v\n", len(body), asString(line, logAsString))
					} else {
						fmt.Printf("           | %v\n", asString(line, logAsString))
					}
					lineNbr++
				} else {
					fmt.Println(asString(line, logAsString))
				}
			}

		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func asString(data []byte, asString bool) interface{} {
	if asString {
		return string(data)
	}
	return data
}
