package main

import (
	"authentication/internals/services"
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Application struct {
	DB  *sql.DB
	svc *services.Services
}

func main() {
	log.Println("Starting Authentication Service")

	// initializing server
	server := initServer()

	// channel to check for server error
	serverError := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			serverError <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("server error: %v", err)
	case err := <-quit:
		log.Printf("Received shutdown signal: %v", err)
	}

	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		return
	}

	log.Println("Server exited properly")
}
