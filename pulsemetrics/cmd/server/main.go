package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/pulsemetrics/pkg/api"
	"github.com/example/pulsemetrics/pkg/ingest"
	"github.com/example/pulsemetrics/pkg/storage"
)

func main() {
	cfg := ingest.DefaultConfig()

	store := storage.NewMemoryStore()
	ing, _ := ingest.NewIngestor(cfg, store)
	defer ing.Close()

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, ing)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("server started on :8080")
		srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("server stopped")
}
