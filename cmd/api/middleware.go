package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/internal/handlers"
	"errors"
	"fmt"
	"net/http"
)

// enableCORS middleware for CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set(
			"Access-Control-Allow-Headers",
			fmt.Sprintf("%s,%s", handlers.AppContentType, handlers.AppJson),
		)
		next.ServeHTTP(w, r)
	})
}

// verifyApiKey middleware to verify if api key provided in the header exists
func verifyApiKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", handlers.AppAuthorization)
		apiKey := r.Header.Get(handlers.AppApiKey)
		if apiKey == "" {
			handlers.Repo.ErrorHandler(w, errors.New("api key not provided"), http.StatusForbidden)
			return
		}
		exists := handlers.Repo.DB.VerifyApiKeyExists(apiKey)
		if !exists {
			handlers.Repo.ErrorHandler(w, errors.New("incorrect API key provided"), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
