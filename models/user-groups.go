package models

import "gorm.io/gorm"

// UserGroup struct for user groups
type UserGroup struct {
	gorm.Model
	Name string `gorm:"type:varchar(50)" json:"name"`
}
