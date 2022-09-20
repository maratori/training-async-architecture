package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/maratori/training-async-architecture/infra"
	"github.com/maratori/training-async-architecture/proto-hub/servicea"
	"github.com/maratori/training-async-architecture/proto-hub/serviceb"
	"github.com/maratori/training-async-architecture/service-b/internal/app"
	"github.com/maratori/training-async-architecture/service-b/internal/domain"
	"github.com/maratori/training-async-architecture/service-b/internal/postgres"
)

func main() {
	db, closeDB, err := infra.NewDB()
	if err != nil {
		panic(err)
	}
	defer closeDB()

	queries := postgres.New(db)
	twirpClient := servicea.NewAServiceJSONClient("service-a", http.DefaultClient)
	service := domain.NewService(app.NewServiceDeps(queries, twirpClient))
	apiService := app.NewBService(service)

	mux := http.NewServeMux()
	twirpServer := serviceb.NewBServiceServer(apiService)
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
