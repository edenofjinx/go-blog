package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/google/uuid"
)

// CreateUser creates a new user (user registration)
func (m *mysqlDatabaseRepo) CreateUser(user models.User) error {
	err := m.DB.Create(&user)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

// UpdateUser updates the user information
func (m *mysqlDatabaseRepo) UpdateUser(user models.User) error {
	err := m.DB.Updates(&user)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

// DeleteUser deletes a user (unregister)
func (m *mysqlDatabaseRepo) DeleteUser(userId int) error {
	result := m.DB.Unscoped().Delete(&models.User{}, userId)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUserByEmail returns the first user with the given email
func (m *mysqlDatabaseRepo) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := m.DB.Model(&models.User{}).Where(&models.User{Email: email}).First(&user)
	if err.Error != nil {
		return user, err.Error
	}
	return user, nil
}

// UpdateUserPassword updates a users password
func (m *mysqlDatabaseRepo) UpdateUserPassword(user models.User) error {
	err := m.DB.Updates(&user)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

// UpdateUserGroup updates a user group
func (m *mysqlDatabaseRepo) UpdateUserGroup(user models.User) error {
	// TODO add func
	return nil
}

// UpdateUserApiKey updates a user api key
func (m *mysqlDatabaseRepo) UpdateUserApiKey(userId int) (string, error) {
	apiKey := uuid.New().String()
	err := m.DB.Model(&models.User{}).Where("id = ?", userId).Update("api_key", apiKey)
	if err.Error != nil {
		return "", err.Error
	}
	return apiKey, nil
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
