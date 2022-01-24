package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetArticlesList get a list of articles with pagination
func (m *mysqlDatabaseRepo) GetArticlesList(c *gin.Context) ([]*models.Article, error) {
	var articles []*models.Article
	rows := m.DB.Scopes(paginate(c.Request, m.App)).Find(&articles)
	if rows.Error != nil {
		m.App.ErrorLog.Println(rows.Error)
		return nil, rows.Error
	}
	return articles, nil
}

// GetArticleById get an article by a given article id
func (m *mysqlDatabaseRepo) GetArticleById(articleId int) (models.ArticleWithContent, error) {
	var article models.ArticleWithContent
	row := m.DB.First(&models.Article{}, articleId)
	row.Scan(&article)
	return article, nil
}

// GetCommentsByArticleId get comments by article id with pagination
func (m *mysqlDatabaseRepo) GetCommentsByArticleId(articleId int, r *http.Request) ([]*models.Comment, error) {
	var comments []*models.Comment
	rows := m.DB.Scopes(paginate(r, m.App)).Where(&models.Comment{ArticleID: articleId}).Find(&comments)
	if rows.Error != nil {
		m.App.ErrorLog.Println(rows.Error)
		return nil, rows.Error
	}
	return comments, nil
}

// VerifyApiKeyExists verify if given api key exists
func (m *mysqlDatabaseRepo) VerifyApiKeyExists(apiKey string) bool {
	var count int64
	if apiKey == "" {
		return false
	}
	m.DB.Model(&models.User{}).Where(&models.User{ApiKey: apiKey}).Count(&count)
	if count == 0 {
		return false
	}
	return true
}

// InsertComment saves comment for given article
func (m *mysqlDatabaseRepo) InsertComment(comment models.Comment) error {
	result := m.DB.Create(&comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *mysqlDatabaseRepo) SaveArticle(article models.Article) error {
	if m.DB.Model(&article).Where("id = ?", article.ID).Updates(&article).RowsAffected == 0 {
		result := m.DB.Create(&article)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
