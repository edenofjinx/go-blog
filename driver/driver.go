package driver

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

// DB hold the database connection pool
type DB struct {
	SQL *gorm.DB
}

var dbConn = &DB{}

// ConnectSQL creates database pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	dbConn.SQL = d
	return dbConn, nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}
