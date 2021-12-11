package models

import (
	"gorm.io/gorm"
)

// MysqlModel is the model for the MySQL database
type MysqlModel struct {
	DB *gorm.DB
}

// AppStatus model for application status response
type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}
