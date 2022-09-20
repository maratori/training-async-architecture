package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/maratori/training-async-architecture/infra"
	"github.com/maratori/training-async-architecture/proto-hub/servicea"
	"github.com/maratori/training-async-architecture/service-a/internal/app"
)

func main() {
	_, closeDB, err := infra.NewDB()
	if err != nil {
		panic(err)
	}
	defer closeDB()

	service := app.NewAService()

	mux := http.NewServeMux()
	twirpServer := servicea.NewAServiceServer(service)
	mux.Handle(twirpServer.PathPrefix(), twirpServer)

	server := http.Server{
		Addr:              ":80",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ErrorLog:          log.Default(),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
