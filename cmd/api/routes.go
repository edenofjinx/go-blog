package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/internal/handlers"
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

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

func unprotectedRoutes(r *httprouter.Router) {
	r.HandlerFunc(http.MethodGet, "/v1/status", handlers.Repo.StatusHandler)
}

func protectedRoutes(r *httprouter.Router, s *alice.Chain) {
	r.GET("/v1/articles/:page/:limit", wrap(s.ThenFunc(handlers.Repo.GetArticlesList)))
	r.GET("/v1/article/:id", wrap(s.ThenFunc(handlers.Repo.GetArticleById)))
	r.GET("/v1/article/:id/comments/:page/:limit", wrap(s.ThenFunc(handlers.Repo.GetCommentsByArticleId)))

	r.POST("/v1/comment/save", wrap(s.ThenFunc(handlers.Repo.SaveComment)))
}
