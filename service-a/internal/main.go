package main

import (
	"context"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/maratori/training-async-architecture/infra"
	"github.com/maratori/training-async-architecture/proto-hub/servicea"
	"github.com/maratori/training-async-architecture/proto-hub/serviceb"
	"github.com/maratori/training-async-architecture/service-a/internal/app"
	"github.com/maratori/training-async-architecture/service-a/internal/domain"
	"github.com/maratori/training-async-architecture/service-a/internal/postgres"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	exitCode := 0
	err := run(ctx)
	if err != nil {
		log.Printf("Exit reason: %+v", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}

func run(ctx context.Context) error {
	db, closeDB, err := infra.NewDB()
	if err != nil {
		return err
	}
	defer closeDB()

	queries := postgres.New(db)
	twirpClient := serviceb.NewBServiceJSONClient("service-b", http.DefaultClient)
	service := domain.NewService(app.NewServiceDeps(queries, twirpClient))
	apiService := app.NewAService(service)

	mux := http.NewServeMux()
	twirpServer := servicea.NewAServiceServer(apiService)
	mux.Handle(twirpServer.PathPrefix(), twirpServer)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return infra.RunHTTPServer(ctx, ":80", mux)
	})

	return eg.Wait()
}
