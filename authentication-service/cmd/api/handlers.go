package main

import (
	"net/http"
)

func (app *Application) authentication(w http.ResponseWriter, r *http.Request) {
	paypal := JsonResponse{
		Error:   false,
		Message: "Broker service is up and running",
	}

	err := app.writeJSON(w, http.StatusOK, paypal)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
