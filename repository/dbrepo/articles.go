package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/gin-gonic/gin"
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

// SaveArticle saves an article
func (m *mysqlDatabaseRepo) SaveArticle(article models.Article) error {
	result := m.DB.Create(&article)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateArticle updates an article
func (m *mysqlDatabaseRepo) UpdateArticle(article models.Article) error {
	result := m.DB.Updates(&article)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteArticle updates an article
func (m *mysqlDatabaseRepo) DeleteArticle(articleId int) error {
	result := m.DB.Delete(&models.Article{}, articleId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
