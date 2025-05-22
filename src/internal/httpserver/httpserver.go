package httpserver

import (
	"github.com/BernsteinMond/brand-scout-test-task/src/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

func New(service service.Service, listenAddr string) *http.Server {
	router := mux.NewRouter()

	server := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	mapHandlers(router, service)

	return server
}
