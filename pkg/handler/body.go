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
)

// LogBody log the body
func LogBody(w http.ResponseWriter, r *http.Request) {
	length := r.Header.Get(headerLogBodyLength) == "true"

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
					fmt.Printf("%10d | %s\n", len(body), line)
				} else {
					fmt.Printf("           | %s\n", line)
				}
				lineNbr++
			} else {
				fmt.Println(line)
			}
		}

	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
