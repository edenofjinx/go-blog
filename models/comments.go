package models

import "gorm.io/gorm"

// Comment struct for comments
type Comment struct {
	gorm.Model
	Content   string  `gorm:"type:text" json:"content"`
	UserID    int     `json:"user_id"`
	User      User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ArticleID int     `json:"article_id"`
	Article   Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
