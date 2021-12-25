package repository

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"net/http"
)

// DatabaseRepository interface for database requests
type DatabaseRepository interface {
	GetArticlesList(r *http.Request) ([]*models.Article, error)
	GetArticleById(r *http.Request) (models.ArticleWithContent, error)
	GetCommentsByArticleId(r *http.Request) ([]*models.Comment, error)
	VerifyApiKeyExists(apiKey string) bool
	InsertComment(comment models.Comment) error
}
