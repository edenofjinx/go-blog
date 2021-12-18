package database

import (
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"errors"
	"gorm.io/gorm"
)

// SeedData seeds database
func SeedData(db *driver.DB) error {
	err := seedUserGroups(db)
	if err != nil {
		return err
	}
	err = seedUsers(db)
	if err != nil {
		return err
	}
	err = seedArticles(db)
	if err != nil {
		return err
	}
	err = seedComments(db)
	if err != nil {
		return err
	}
	return nil
}

// seedUserGroups seeds user groups if empty
func seedUserGroups(db *driver.DB) error {
	var userGroups = []models.UserGroup{
		{
			Name: "Admin",
		},
		{
			Name: "Registered",
		},
	}
	if err := db.SQL.AutoMigrate(&models.UserGroup{}); err == nil && db.SQL.Migrator().HasTable(&models.UserGroup{}) {
		if err := db.SQL.First(&models.UserGroup{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			for _, v := range userGroups {
				db.SQL.Create(&v)
			}
		}
	}
	return nil
}

// seedUsers seeds users table if empty
func seedUsers(db *driver.DB) error {
	var users = []models.User{
		{
			Name:    "Admin",
			ApiKey:  "8c3751f5-39f8-4672-8362-1d83e3169ae3",
			GroupID: 1,
		},
		{
			Name:    "Registered",
			ApiKey:  "80f45b24-874b-4e96-9e69-0efd000eca4a",
			GroupID: 2,
		},
	}
	if err := db.SQL.AutoMigrate(&models.User{}); err == nil && db.SQL.Migrator().HasTable(&models.User{}) {
		if err := db.SQL.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			for _, v := range users {
				db.SQL.Create(&v)
			}
		}
	}
	return nil
}

// seedArticles seeds articles if empty
func seedArticles(db *driver.DB) error {
	var articles = []models.Article{
		{
			Title:   "Test 1",
			Content: "Some test content 1",
			UserID:  1,
		},
		{
			Title:   "Test 2",
			Content: "Test content 2",
			UserID:  2,
		},
	}
	if err := db.SQL.AutoMigrate(&models.Article{}); err == nil && db.SQL.Migrator().HasTable(&models.Article{}) {
		if err := db.SQL.First(&models.Article{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			for _, v := range articles {
				db.SQL.Create(&v)
			}
		}
	}
	return nil
}

// seedComments seeds comments if empty
func seedComments(db *driver.DB) error {
	var comments = []models.Comment{
		{
			ArticleID: 1,
			Content:   "test comment 1",
			UserID:    1,
		},
		{
			ArticleID: 2,
			Content:   "Test comment 2",
			UserID:    2,
		},
	}
	if err := db.SQL.AutoMigrate(&models.Comment{}); err == nil && db.SQL.Migrator().HasTable(&models.Comment{}) {
		if err := db.SQL.First(&models.Comment{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			for _, v := range comments {
				db.SQL.Create(&v)
			}
		}
	}
	return nil
}
