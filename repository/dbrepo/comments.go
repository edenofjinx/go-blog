package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"net/http"
)

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

// SaveComment saves comment for given article
func (m *mysqlDatabaseRepo) SaveComment(comment models.Comment) error {
	result := m.DB.Create(&comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateComment updates a comment
func (m *mysqlDatabaseRepo) UpdateComment(comment models.Comment) error {
	result := m.DB.Updates(&comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteComment deletes a comment by id
func (m *mysqlDatabaseRepo) DeleteComment(commentID int) error {
	result := m.DB.Delete(&models.Comment{}, commentID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
