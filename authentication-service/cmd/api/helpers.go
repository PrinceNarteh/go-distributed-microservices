package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != nil {
		return errors.New("body must have only a single json value")
	}

	return nil
}
func (app *Application) writeJSON(w http.ResponseWriter, statusCode int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Application) errorJSON(w http.ResponseWriter, err error, statusCode ...int) error {
	status := http.StatusBadRequest

	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, status, payload)
}
