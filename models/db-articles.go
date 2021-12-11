package models

import (
	"gorm.io/gorm"
)

// NewDatabase returns model with the database pool
func NewDatabase(db *gorm.DB) MysqlModel {
	return MysqlModel{
		DB: db,
	}
}

func SeedData(db *gorm.DB) error {
	//if err := db.AutoMigrate(&Product{}); err == nil && db.Migrator().HasTable(&Product{}) {
	//	if err := db.First(&Product{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
	//		db.Create(&Product{Code: "D42", Price: 100})
	//	}
	//}
	//if err != nil {
	//	log.Fatal(err)
	//	return err
	//}
	return nil
}
