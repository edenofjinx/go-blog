package models

import (
	"gorm.io/gorm"
)

// Article struct for articles
type Article struct {
	gorm.Model
	Title   string `gorm:"type:varchar(100)" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	UserID  int    `json:"user_id"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
