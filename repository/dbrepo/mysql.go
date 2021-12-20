package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (m *mysqlDatabaseRepo) GetArticlesList(r *http.Request) ([]*models.Article, error) {
	var articles []*models.Article
	rows := m.DB.Scopes(paginate(r, m.App)).Find(&articles)
	if rows.Error != nil {
		m.App.ErrorLog.Println(rows.Error)
		return nil, rows.Error
	}
	return articles, nil
}

func (m *mysqlDatabaseRepo) GetArticleById(r *http.Request) (models.ArticleWithContent, error) {
	var article models.ArticleWithContent
	params := httprouter.ParamsFromContext(r.Context())
	articleId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		m.App.ErrorLog.Println(err)
		return article, err
	}
	row := m.DB.First(&models.Article{}, articleId)
	row.Scan(&article)
	return article, nil
}
func (m *mysqlDatabaseRepo) GetCommentsByArticleId(r *http.Request) ([]*models.Comment, error) {
	var comments []*models.Comment
	params := httprouter.ParamsFromContext(r.Context())
	articleId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		m.App.ErrorLog.Println(err)
		return nil, err
	}
	rows := m.DB.Scopes(paginate(r, m.App)).Where(&models.Comment{ArticleID: articleId}).Find(&comments)
	if rows.Error != nil {
		m.App.ErrorLog.Println(rows.Error)
		return nil, rows.Error
	}
	return comments, nil
}

func (m *mysqlDatabaseRepo) VerifyApiKeyExists(apiKey string) bool {
	var count int64
	m.DB.Model(&models.User{}).Where(&models.User{ApiKey: apiKey}).Count(&count)
	if count == 0 {
		return false
	}
	return true
}

func paginate(r *http.Request, config *config.AppConfig) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		params := httprouter.ParamsFromContext(r.Context())
		page, err := strconv.Atoi(params.ByName("page"))
		if err != nil {
			config.ErrorLog.Println(err)
		}
		if page == 0 {
			page = 1
		}

		limit, err := strconv.Atoi(params.ByName("limit"))
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
