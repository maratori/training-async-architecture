package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	server := http.Server{
		Addr: ":80",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			defer r.Body.Close()
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(body)
		}),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ErrorLog:          log.Default(),
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
