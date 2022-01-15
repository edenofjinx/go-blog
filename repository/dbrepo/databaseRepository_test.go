package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

var a *config.AppConfig
var db *gorm.DB

func TestNewMysqlRepo(t *testing.T) {
	repo := NewMysqlRepo(db, a)
	assert.ObjectsAreEqualValues(mysqlDatabaseRepo{}, repo)
}
