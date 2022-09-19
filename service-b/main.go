package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe("localhost:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer r.Body.Close()
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	if err != nil {
		log.Fatalln(err)
	}
}
