package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"personal_erp/backend/internal/api"
	"personal_erp/backend/internal/config"
	"personal_erp/backend/internal/db"
	"personal_erp/backend/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(pool); err != nil {
		log.Fatalf("migrate database: %v", err)
	}
	log.Printf("database ready (%s)", cfg.DBName)

	server := api.NewServer(pool, store.NewActivityStore(pool), store.NewCategoryStore(pool))

	httpServer := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: server.Routes(),
	}

	go func() {
		log.Printf("Personal ERP API listening on http://localhost:%s", cfg.ServerPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	log.Println("server stopped")
}
