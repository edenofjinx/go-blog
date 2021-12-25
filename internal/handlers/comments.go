package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"net/http"
	"time"
)

// GetCommentsByArticleId handler to get comments by article id
func (repo *Repository) GetCommentsByArticleId(w http.ResponseWriter, r *http.Request) {
	comments, err := repo.DB.GetCommentsByArticleId(r)
	if err != nil {
		repo.App.ErrorLog.Println(err)
		return
	}
	js, err := json.MarshalIndent(comments, "", "\t")
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

// SaveComment saves comment into the database
func (repo *Repository) SaveComment(w http.ResponseWriter, r *http.Request) {
	var payload models.CommentPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		repo.ErrorHandler(w, err)
		return
	}
	content, err := repo.parseImageTags(payload.Content)
	if err != nil {
		repo.ErrorHandler(w, err)
		return
	}
	var comment models.Comment
	comment.Content = content
	comment.UserID = payload.UserID
	comment.ArticleID = payload.ArticleID
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	err = repo.DB.InsertComment(comment)
	if err != nil {
		repo.ErrorHandler(w, err)
		return
	}
	message := JsonResponse{
		Message: "Comment has been successfully saved.",
	}
	err = repo.ResponseJson(w, http.StatusOK, message, "success")
	if err != nil {
		repo.ErrorHandler(w, err)
		return
	}
}
