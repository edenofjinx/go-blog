package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/repository"
	"bitbucket.org/julius_liaudanskis/go-blog/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepository
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewMysqlRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}
