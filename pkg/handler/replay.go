package handler

import (
	"io"
	"net/http"
)

const (
	headerTrainReplay = "Train-Replay"
)

var (
	replayBody        = map[string][]byte{}
	replayContentType = map[string]string{}
)

// Replay replay
func Replay(w http.ResponseWriter, r *http.Request) {
	train := r.Header.Get(headerTrainReplay)

	if train == "true" {
		body, err := io.ReadAll(r.Body)
		if err == nil {
			replayBody[r.RequestURI] = body
		}
		replayContentType[r.RequestURI] = r.Header.Get("Content-Type")
		defer func() { _ = r.Body.Close() }()
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
