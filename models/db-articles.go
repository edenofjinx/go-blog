package models

import (
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
)

func SeedData(db *driver.DB) error {
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
