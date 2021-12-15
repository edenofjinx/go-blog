package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/internal/handlers"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// routes holds data of all available routes for the app
func routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", handlers.Repo.StatusHandler)
	return router
}
