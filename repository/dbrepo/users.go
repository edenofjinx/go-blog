package dbrepo

import "bitbucket.org/julius_liaudanskis/go-blog/models"

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
