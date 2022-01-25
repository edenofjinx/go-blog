package dbrepo

import "bitbucket.org/julius_liaudanskis/go-blog/models"

func (m *mysqlDatabaseRepo) CreateUser(user models.User) error {
	err := m.DB.Create(&user)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (m *mysqlDatabaseRepo) UpdateUser(user models.User) error {
	err := m.DB.Updates(&user)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (m *mysqlDatabaseRepo) DeleteUser(userId int) error {
	result := m.DB.Unscoped().Delete(&models.User{}, userId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *mysqlDatabaseRepo) UpdateUserPassword(user models.User) error {
	return nil
}

func (m *mysqlDatabaseRepo) UpdateUserGroup(user models.User) error {
	return nil
}

func (m *mysqlDatabaseRepo) UpdateUserApiKey(userId int) error {
	return nil
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
