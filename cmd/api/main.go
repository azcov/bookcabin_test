package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azcov/bookcabin_test/internal/config"
	"github.com/azcov/bookcabin_test/internal/service"
	"github.com/azcov/bookcabin_test/internal/transport/api"
	"github.com/azcov/bookcabin_test/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Build dependencies

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found or error loading it", "err", err)
	}

	cfg := config.NewConfig()
	config.LoadConfig(cfg)
	logger.Info("Loading config", "cfg", cfg)
	svc := service.NewFlightService(*cfg)
	h := api.NewHandler(svc)
	r := api.NewRouter(h)

	// Start http.Server and graceful shutdown
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
