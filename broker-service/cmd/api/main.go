package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "8080"

type Application struct{}

func main() {
	app := Application{}

	log.Println("starting broker-service on port:", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
