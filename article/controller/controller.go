package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}


// Register all Controllers and its Routes
func New(logger *log.Logger) *mux.Router {

	// Initializing Article Handler
	articleHandle := newArticle(logger)

	// create a new serve mux and register handlers
	sm := mux.NewRouter()

	// Handle All GET
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/article/", responseHandler(articleHandle.GetAll))
	getRouter.HandleFunc("/article", responseHandler(articleHandle.GetAll))
	getRouter.HandleFunc("/article/{id:[0-9)]+}", responseHandler(articleHandle.Get))
	getRouter.HandleFunc("/article/{id:[0-9)]+}/", responseHandler(articleHandle.Get))

	// Handle All PUT
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/article/{id:[0-9)]+}", responseHandler(articleHandle.Put))
	putRouter.HandleFunc("/article/{id:[0-9)]+}/", responseHandler(articleHandle.Put))

	// Handle All POST
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/article/", responseHandler(articleHandle.Create))
	postRouter.HandleFunc("/article", responseHandler(articleHandle.Create))

	// Handle All DELETE
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/article/{id:[0-9)]+}", responseHandler(articleHandle.Delete))
	deleteRouter.HandleFunc("/article/{id:[0-9)]+}/", responseHandler(articleHandle.Delete))

	return sm
}

// responseHandler format response into json and also handle error
func responseHandler(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		// Add Cors
		wr.Header().Set("Access-Control-Allow-Origin", "*")
		data, status, err := h(wr, req)
		if err != nil {
			data = err.Error()
		}
		wr.Header().Set("Content-Type", "application/json")
		wr.WriteHeader(status)
		err = json.NewEncoder(wr).Encode(response{Data: data, Success: err == nil})
		if err != nil {
			log.Printf("could not encode response to output: %v", err)
		}

	}
}

