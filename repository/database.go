package repository

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DatabaseRepository interface for database requests
type DatabaseRepository interface {
	articleRepository
	commentRepository
	userRepository
}

type articleRepository interface {
	GetArticlesList(c *gin.Context) ([]*models.Article, error)
	GetArticleById(articleId int) (models.ArticleWithContent, error)
	SaveArticle(article models.Article) error
	UpdateArticle(article models.Article) error
	DeleteArticle(articleId int) error
}

type commentRepository interface {
	GetCommentsByArticleId(articleId int, r *http.Request) ([]*models.Comment, error)
	SaveComment(comment models.Comment) error
	UpdateComment(comment models.Comment) error
	DeleteComment(commentId int) error
}

type userRepository interface {
	CreateUser(user models.User) error
	UpdateUser(user models.User) error
	DeleteUser(userId int) error
	GetUserByEmail(email string) (models.User, error)
	UpdateUserPassword(user models.User) error
	UpdateUserGroup(user models.User) error
	UpdateUserApiKey(userId int) error
	VerifyApiKeyExists(apiKey string) bool
}

type userGroupRepository interface {
	//TODO
}
