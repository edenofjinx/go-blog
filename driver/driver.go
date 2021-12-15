package driver

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB hold the database connection pool
type DB struct {
	SQL *gorm.DB
}

var dbConn = &DB{}

// ConnectSQL creates database pool for Postgres
func ConnectSQL(dsn string, app config.AppConfig) (*DB, error) {
	d, err := NewDatabase(dsn, app)
	if err != nil {
		panic(err)
	}
	dbConn.SQL = d
	return dbConn, nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string, app config.AppConfig) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}
	return db, nil
}
