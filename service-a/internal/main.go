package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/maratori/training-async-architecture/infra"
	"github.com/maratori/training-async-architecture/proto-hub/servicea"
	"github.com/maratori/training-async-architecture/proto-hub/serviceb"
	"github.com/maratori/training-async-architecture/service-a/internal/app"
	"github.com/maratori/training-async-architecture/service-a/internal/domain"
	"github.com/maratori/training-async-architecture/service-a/internal/postgres"
)

func main() {
	db, closeDB, err := infra.NewDB()
	if err != nil {
		panic(err)
	}
	defer closeDB()

	queries := postgres.New(db)
	twirpClient := serviceb.NewBServiceJSONClient("service-b", http.DefaultClient)
	service := domain.NewService(app.NewServiceDeps(queries, twirpClient))
	apiService := app.NewAService(service)

	mux := http.NewServeMux()
	twirpServer := servicea.NewAServiceServer(apiService)
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
