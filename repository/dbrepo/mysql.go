package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func (m *mysqlDatabaseRepo) AllStatus() bool {
	return true
}

func (m *mysqlDatabaseRepo) GetArticlesList(r *http.Request) ([]*models.Article, error) {
	var articles []*models.Article
	rows := m.DB.Scopes(paginate(r)).Find(&articles)
	if rows.Error != nil {
		log.Println(rows.Error)
		return nil, rows.Error
	}
	return articles, nil
}

func (m *mysqlDatabaseRepo) GetArticleById(r *http.Request) (models.ArticleWithContent, error) {
	var article models.ArticleWithContent
	params := httprouter.ParamsFromContext(r.Context())
	articleId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return nil, err
	}
	rows := m.DB.Scopes(paginate(r)).Where(&models.Comment{ArticleID: articleId}).Find(&comments)
	if rows.Error != nil {
		log.Println(rows.Error)
		return nil, rows.Error
	}
	return comments, nil
}

func paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		params := httprouter.ParamsFromContext(r.Context())
		page, err := strconv.Atoi(params.ByName("page"))
		if err != nil {
			log.Println(err)
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
