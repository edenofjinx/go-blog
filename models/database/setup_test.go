package database

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
)

var testDsn = "root:root@tcp(localhost:3306)/go_blog_test?charset=utf8&parseTime=True&loc=Local"

type seedTestSuite struct {
	suite.Suite
	db *driver.DB
}

func (suite *seedTestSuite) SetupSuite() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	a := config.AppConfig{
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		InProduction: false,
		AppVersion:   "test",
		Environment:  "test",
		StaticImages: "static/test/images/",
	}
	d, err := driver.ConnectSQL(testDsn, a)
	if err != nil {
		suite.Fail("could not connect to the test sql")
	}
	suite.db = d
}

func (suite *seedTestSuite) TearDownTest() {
	err := suite.db.SQL.Migrator().DropTable(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not drop test tables")
}

func (suite *seedTestSuite) TearDownSuite() {
	err := suite.db.SQL.Migrator().DropTable(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not drop test tables")
}
