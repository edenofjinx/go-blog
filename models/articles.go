package models

import (
	"gorm.io/gorm"
)

// Article struct for articles
type Article struct {
	gorm.Model
	Title   string `gorm:"type:varchar(100)" json:"title"`
	Content string `gorm:"type:text" json:"-"`
	UserID  int    `json:"-"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

// ArticleWithContent returns Article but with content in the json
type ArticleWithContent struct {
	Article
	Content string `gorm:"type:text" json:"content"`
}
