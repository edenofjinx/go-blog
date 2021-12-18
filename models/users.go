package models

import "gorm.io/gorm"

// User struct for users
type User struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100)" json:"name"`
	ApiKey  string `gorm:"type:varchar(60)" json:"api_key"`
	GroupID int
	Group   UserGroup `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
