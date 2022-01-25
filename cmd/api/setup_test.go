package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

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

type TestMainPackage struct {
	suite.Suite
}

func (suite *TestMainPackage) SetupSuite() {

}

func (suite *TestMainPackage) SetupTest() {
	os.Clearenv()
}

func (suite *TestMainPackage) TearDownSuite() {
	os.Clearenv()
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(TestMainPackage))
}
