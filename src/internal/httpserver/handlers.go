package httpserver

import (
	"encoding/json"
	"errors"
	quoteService "github.com/BernsteinMond/brand-scout-test-task/src/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func mapHandlers(router *mux.Router, service quoteService.Service) {
	quotesGroup := router.PathPrefix("/quotes").Subrouter()
	quotesGroup.Handle("", postQuoteHandler(service)).Methods("POST")
	quotesGroup.Handle("", getQuotesHandler(service)).Methods("GET")
	quotesGroup.Handle("/random", getRandomQuoteHandler(service)).Methods("GET")
	quotesGroup.Handle("/{id}", deleteQuoteHandler(service)).Methods("DELETE")
}

func postQuoteHandler(service quoteService.Service) http.HandlerFunc {
	type request = quoteCreateDTO
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}

		if len(req.Quote) == 0 || len(req.Author) == 0 {
			http.Error(w, "empty \"quote\" or \"author\" parameter", http.StatusBadRequest)
		}

		err = service.CreateNewQuote(r.Context(), req.Author, req.Quote)
		if err != nil {
			if errors.Is(err, quoteService.ErrAlreadyExists) {
				w.WriteHeader(http.StatusConflict)
				return
			}

			http.Error(w, "service: create new quote", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func getQuotesHandler(service quoteService.Service) http.HandlerFunc {
	type response struct {
		Quotes []quoteReadDTO `json:"quotes"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		authorFilter := r.URL.Query().Get("author")

		quotes, err := service.GetQuotesWithFilter(r.Context(), authorFilter)
		if err != nil {
			http.Error(w, "service: get quotes", http.StatusInternalServerError)
			return
		}

		resp := response{
			Quotes: make([]quoteReadDTO, len(quotes)),
		}

		for i, quote := range quotes {
			resp.Quotes[i] = quoteFromDomainToReadDTO(&quote)
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func getRandomQuoteHandler(service quoteService.Service) http.HandlerFunc {
	type response = quoteReadDTO
	return func(w http.ResponseWriter, r *http.Request) {
		quote, err := service.GetRandomQuote(r.Context())
		if err != nil {
			http.Error(w, "service: get random quote", http.StatusInternalServerError)
			return
		}

		var resp response
		resp = quoteFromDomainToReadDTO(quote)

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}

		return
	}
}

func deleteQuoteHandler(service quoteService.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, ok := mux.Vars(r)["id"]
		if !ok {
			http.Error(w, "empty \"id\" parameter", http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid \"id\" parameter", http.StatusBadRequest)
			return
		}

		err = service.DeleteQuoteByID(r.Context(), id)
		if err != nil {
			http.Error(w, "service: delete quote by id", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
