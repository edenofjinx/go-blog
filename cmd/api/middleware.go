package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/internal/handlers"
	"errors"
	"fmt"
	"net/http"
)

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

func checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", handlers.AppAuthorization)
		apiKey := r.Header.Get(handlers.AppApiKey)
		if apiKey == "" {
			handlers.Repo.ErrorHandler(w, errors.New("invalid auth header"), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
