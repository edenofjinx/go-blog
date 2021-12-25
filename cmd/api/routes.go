package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/internal/handlers"
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

// wrap wraps the handler to return httprouter
func wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// routes holds data of all available routes for the app
func routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(verifyApiKey)
	router.ServeFiles("/static/images/*filepath", http.Dir("static/images"))
	unprotectedRoutes(router)

	protectedRoutes(router, &secure)
	return enableCORS(router)
}

// unprotectedRoutes holds routes that are not protected by an api key
func unprotectedRoutes(r *httprouter.Router) {
	r.GET("/v1/status", handlers.Repo.StatusHandler)
}

// protectedRoutes holds routes that are protected by an api key
func protectedRoutes(r *httprouter.Router, s *alice.Chain) {
	r.GET("/v1/articles", wrap(s.ThenFunc(handlers.Repo.GetArticlesList)))
	r.GET("/v1/article/:id", wrap(s.ThenFunc(handlers.Repo.GetArticleById)))
	r.GET("/v1/article/:id/comments", wrap(s.ThenFunc(handlers.Repo.GetCommentsByArticleId)))

	r.POST("/v1/comment/save", wrap(s.ThenFunc(handlers.Repo.SaveComment)))
}
