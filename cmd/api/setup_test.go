package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
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

type TestMainPackage struct {
	suite.Suite
}

func (suite *TestMainPackage) SetupSuite() {
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
	log.Println(a)
}

func (suite *TestMainPackage) TearDownSuite() {

}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestMainSuite(t *testing.T) {
	suite.Run(t, new(TestMainPackage))
}

type testPagination struct {
	limit string
	page  string
	order string
}

func generateNewGETRequest(url string, pagination testPagination) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if pagination.limit != "" {
		q := req.URL.Query()
		q.Add("limit", pagination.limit)
		req.URL.RawQuery = q.Encode()
	}
	if pagination.page != "" {
		q := req.URL.Query()
		q.Add("page", pagination.page)
		req.URL.RawQuery = q.Encode()
	}
	if pagination.order != "" {
		q := req.URL.Query()
		q.Add("order", pagination.order)
		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}
