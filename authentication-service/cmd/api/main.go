package main

import (
	"authentication/internals/config"
	"authentication/internals/db"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Application struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting Authentication Service")

	// connect database
	db := db.ConnectDB()
	if db == nil {
		log.Fatal("could not connect to Postgres")
	}

	// initializing app
	app := Application{
		DB: db,
	}

	webPort := config.Env.App.Port
	log.Println("starting authentication-service on port:", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
