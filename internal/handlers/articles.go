package handlers

import (
	"net/http"
)

// GetArticlesList handler to get articles list
func (repo *Repository) GetArticlesList(w http.ResponseWriter, r *http.Request) {
	articles, err := repo.DB.GetArticlesList(r)
	if err != nil {
		repo.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	err = repo.ResponseJson(w, http.StatusAccepted, articles, "data")
	if err != nil {
		repo.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
}

// GetArticleById handler to get article data by article id
func (repo *Repository) GetArticleById(w http.ResponseWriter, r *http.Request) {
	articles, err := repo.DB.GetArticleById(r)
	if err != nil {
		repo.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	err = repo.ResponseJson(w, http.StatusAccepted, articles, "data")
	if err != nil {
		repo.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
}
