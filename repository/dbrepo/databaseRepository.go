package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/repository"
	"gorm.io/gorm"
)

type mysqlDatabaseRepo struct {
	App *config.AppConfig
	DB  *gorm.DB
}

func NewMysqlRepo(conn *gorm.DB, a *config.AppConfig) repository.DatabaseRepository {
	return &mysqlDatabaseRepo{
		App: a,
		DB:  conn,
	}
}
