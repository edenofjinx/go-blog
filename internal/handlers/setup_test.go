package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
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

type testPagination struct {
	limit string
	page  string
	order string
}

type handlersTestSuite struct {
	suite.Suite
	testHandlerRepo *Repository
	testSQL         *driver.DB
}

func (suite *handlersTestSuite) SetupSuite() {
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
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
	err := godotenv.Load(root + "/.env.test")
	suite.Nil(err, "should not throw an error when loading .env file")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_URL"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE"),
	)
	db, err := driver.ConnectSQL(dsn, a)
	if err != nil {
		suite.Fail("could not connect to the test sql")
	}
	handlerRepo := NewRepo(&a, db)
	NewHandlers(handlerRepo)
	suite.testHandlerRepo = handlerRepo
	suite.testSQL = db
	err = suite.testSQL.SQL.AutoMigrate(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not create test tables")
	suite.testSQL.SQL.Create(&testUserGroups)
	suite.testSQL.SQL.Create(&testUsers)
	suite.testSQL.SQL.Create(&testArticles)
	suite.testSQL.SQL.Create(&testComments)
}

func (suite *handlersTestSuite) TearDownSuite() {
	err := suite.testSQL.SQL.Migrator().DropTable(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	suite.Nil(err, "could not drop test tables")
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../../")
	err = os.RemoveAll(root + "/" + suite.testHandlerRepo.App.StaticImages)
	suite.Nil(err, "could not remove test images from folder")
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(handlersTestSuite))
}

func generateNewGETRequest(testUrl string, pagination testPagination) (*gin.Context, *httptest.ResponseRecorder) {
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	req := &http.Request{
		URL: &url.URL{
			Path: testUrl,
		},
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
	c.Request = req
	return c, rr
}
