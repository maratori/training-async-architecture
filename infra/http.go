package infra

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func RunHTTPServer(ctx context.Context, addr string, handler http.Handler) error {
	const (
		readHeaderTimeout = 10 * time.Second
		writeTimeout      = 10 * time.Second
		shutdownTimeout   = 1 * time.Second
	)

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		ErrorLog:          log.Default(),
	}

	serverErrsChannel := make(chan error, 1)
	go func() {
		serverErrsChannel <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err := srv.Shutdown(ctxShutdown) //nolint:contextcheck // need to use background context because ctx is closed.
		if err != nil {
			return fmt.Errorf("http.Server.Shutdown(%s): %w", addr, err)
		}
		return nil
	case err := <-serverErrsChannel:
		if err != nil {
			return fmt.Errorf("http.Server.ListenAndServe(%s): %w", addr, err)
		}
		return nil
	}
}
