package main

import (
	"authentication/internals/config"
	"authentication/internals/db"
	"authentication/internals/repositories"
	"authentication/internals/services"
	"log"
	"net/http"
)

func initServer() *http.Server {
	db := db.ConnectDB()
	if db != nil {
		log.Fatal("could not connect to Postgres")
	}
	repo := repositories.NewRepository(db)
	services := services.NewService(repo)

	// initializing app
	app := Application{
		DB:  db,
		svc: services,
	}

	webPort := config.Env.App.Port
	log.Println("Authentication Service started on port", webPort)

	server := &http.Server{
		Addr:    webPort,
		Handler: app.routes(),
	}

	return server
}
