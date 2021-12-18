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
	router.HandlerFunc(http.MethodGet, "/articles/:page/:limit", handlers.Repo.GetArticlesList)
	router.HandlerFunc(http.MethodGet, "/article/:id", handlers.Repo.GetArticleById)
	router.HandlerFunc(http.MethodGet, "/article/:id/comments/:page/:limit", handlers.Repo.GetCommentsByArticleId)
	return router
}
