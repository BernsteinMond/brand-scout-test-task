package httpserver

import (
	"github.com/gorilla/mux"
	"net/http"
)

func New(service QuoteService, listenAddr string) *http.Server {
	router := mux.NewRouter()

	server := &http.Server{
		Addr:    ":" + listenAddr,
		Handler: router,
	}

	mapHandlers(router, service)

	return server
}
