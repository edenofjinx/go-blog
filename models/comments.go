package models

import "gorm.io/gorm"

// Comment struct for comments
type Comment struct {
	gorm.Model
	Content   string  `gorm:"type:text" json:"content"`
	UserID    int     `json:"user_id"`
	User      User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	ArticleID int     `json:"article_id"`
	Article   Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

// CommentPayload struct for comment payload
type CommentPayload struct {
	Content   string `json:"content"`
	UserID    int    `json:"user_id"`
	ArticleID int    `json:"article_id"`
}
