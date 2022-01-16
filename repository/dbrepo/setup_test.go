package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

var testDsn = "root:root@tcp(localhost:3306)/go_blog_test?charset=utf8&parseTime=True&loc=Local"

var testUserGroups = []models.UserGroup{
	{
		Name: "Test Admin",
	},
	{
		Name: "Test Registered",
	},
}

var testUsers = []models.User{
	{
		Name:    "Test Admin",
		ApiKey:  "11111111-2222-3333-4444-5555555555555",
		GroupID: 1,
	},
	{
		Name:    "Test Registered",
		ApiKey:  "66666666-7777-8888-9999-101010101010",
		GroupID: 2,
	},
}

var testArticles = []models.Article{
	{
		Title:   "Test 1",
		Content: "Some test content 1",
		UserID:  1,
	},
	{
		Title:   "Test 2",
		Content: "Test content 2",
		UserID:  2,
	},
}

var testComments = []models.Comment{
	{
		ArticleID: 1,
		Content:   "test comment 1",
		UserID:    1,
	},
	{
		ArticleID: 1,
		Content:   "test comment 1-2",
		UserID:    2,
	},
	{
		ArticleID: 2,
		Content:   "Test comment 2",
		UserID:    2,
	},
	{
		ArticleID: 2,
		Content:   "Test comment 2-1",
		UserID:    1,
	},
}

type testPagination struct {
	limit string
	page  string
	order string
}

type databaseRequestTestSuite struct {
	suite.Suite
	testRepo mysqlDatabaseRepo
}

func (suite *databaseRequestTestSuite) SetupSuite() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	a := config.AppConfig{
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		InProduction: false,
		AppVersion:   "test",
		Environment:  "test",
		StaticImages: "static/test/",
	}
	db, err := driver.ConnectSQL(testDsn, a)
	if err != nil {
		suite.Fail("could not connect to the test sql")
	}
	suite.testRepo.App = &a
	suite.testRepo.DB = db.SQL
	err = suite.testRepo.DB.AutoMigrate(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not create test tables")
	suite.testRepo.DB.Create(&testUserGroups)
	suite.testRepo.DB.Create(&testUsers)
	suite.testRepo.DB.Create(&testArticles)
	suite.testRepo.DB.Create(&testComments)
}

func (suite *databaseRequestTestSuite) TearDownSuite() {
	err := suite.testRepo.DB.Migrator().DropTable(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not drop test tables")
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(databaseRequestTestSuite))
}
