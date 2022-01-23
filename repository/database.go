package repository

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/gin-gonic/gin"
)

// DatabaseRepository interface for database requests
type DatabaseRepository interface {
	GetArticlesList(c *gin.Context) ([]*models.Article, error)
	GetArticleById(c *gin.Context) (models.ArticleWithContent, error)
	GetCommentsByArticleId(c *gin.Context) ([]*models.Comment, error)
	VerifyApiKeyExists(apiKey string) bool
	InsertComment(comment models.Comment) error
}
