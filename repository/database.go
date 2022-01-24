package repository

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DatabaseRepository interface for database requests
type DatabaseRepository interface {
	GetArticlesList(c *gin.Context) ([]*models.Article, error)
	GetArticleById(articleId int) (models.ArticleWithContent, error)
	GetCommentsByArticleId(articleId int, r *http.Request) ([]*models.Comment, error)
	VerifyApiKeyExists(apiKey string) bool
	InsertComment(comment models.Comment) error
	SaveArticle(article models.Article) error
}
