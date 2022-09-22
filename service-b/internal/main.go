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
	"github.com/maratori/training-async-architecture/service-b/internal/app"
	"github.com/maratori/training-async-architecture/service-b/internal/domain"
	"github.com/maratori/training-async-architecture/service-b/internal/postgres"
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
	twirpClient := servicea.NewAServiceJSONClient("service-a", http.DefaultClient)
	service := domain.NewService(app.NewServiceDeps(queries, twirpClient))
	apiService := app.NewBService(service)

	mux := http.NewServeMux()
	twirpServer := serviceb.NewBServiceServer(apiService)
	mux.Handle(twirpServer.PathPrefix(), twirpServer)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return infra.RunHTTPServer(ctx, ":80", mux)
	})

	return eg.Wait()
}
