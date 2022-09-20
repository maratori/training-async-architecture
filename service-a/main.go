package main

import (
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/maratori/training-async-architecture/infra"
)

func main() {
	_, closeDB, err := infra.NewDB()
	if err != nil {
		panic(err)
	}
	defer closeDB()

	server := http.Server{
		Addr: ":80",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, errR := io.ReadAll(r.Body)
			if errR != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(errR.Error()))
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
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
