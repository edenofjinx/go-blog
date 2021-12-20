package handlers

import (
	"encoding/json"
	"net/http"
)

// GetArticlesList handler to get articles list
func (repo *Repository) GetArticlesList(w http.ResponseWriter, r *http.Request) {
	articles, err := repo.DB.GetArticlesList(r)
	if err != nil {
		repo.App.ErrorLog.Println(err)
		return
	}
	js, err := json.MarshalIndent(articles, "", "\t")
	if err != nil {
		repo.App.ErrorLog.Println(err)
	}
	w.Header().Set(AppContentType, AppJson)
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write(js)
	if err != nil {
		repo.App.ErrorLog.Println(err)
		return
	}
}

// GetArticleById handler to get article data by article id
func (repo *Repository) GetArticleById(w http.ResponseWriter, r *http.Request) {
	articles, err := repo.DB.GetArticleById(r)
	if err != nil {
		repo.App.ErrorLog.Println(err)
		return
	}
	js, err := json.MarshalIndent(articles, "", "\t")
	if err != nil {
		repo.App.ErrorLog.Println(err)
	}
	w.Header().Set(AppContentType, AppJson)
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write(js)
	if err != nil {
		repo.App.ErrorLog.Println(err)
		return
	}
}
