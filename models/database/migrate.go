package database

import (
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
)

func MigrateData(db *driver.DB) error {
	err := db.SQL.AutoMigrate(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	if err != nil {
		return err
	}
	return nil
}
