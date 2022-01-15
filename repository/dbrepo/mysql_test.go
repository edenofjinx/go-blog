package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/config"
	"bitbucket.org/julius_liaudanskis/go-blog/driver"
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
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

var testsForGetArticlesList = []struct {
	name            string
	pagination      testPagination
	expectedLength  int
	expectedTitle   string
	expectedContent string
}{
	{
		name:            "default articles list",
		pagination:      testPagination{limit: "", page: "", order: ""},
		expectedLength:  2,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
	{
		name:            "articles list with limit of 1",
		pagination:      testPagination{limit: "1", page: "", order: ""},
		expectedLength:  1,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
	{
		name:            "articles list with limit of 1000",
		pagination:      testPagination{limit: "1000", page: "", order: ""},
		expectedLength:  2,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
	{
		name:            "articles list with page of 5",
		pagination:      testPagination{limit: "", page: "5", order: ""},
		expectedLength:  0,
		expectedTitle:   "",
		expectedContent: "",
	},
	{
		name:            "articles list with order of DESC",
		pagination:      testPagination{limit: "", page: "", order: "DESC"},
		expectedLength:  2,
		expectedTitle:   testArticles[1].Title,
		expectedContent: testArticles[1].Content,
	},
	{
		name:            "articles list with mixed pagination",
		pagination:      testPagination{limit: "1", page: "2", order: "DESC"},
		expectedLength:  1,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
	{
		name:            "articles list with incorrect pagination variables",
		pagination:      testPagination{limit: "test", page: "-50", order: "TestOrder"},
		expectedLength:  2,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
	},
}

var testsForGetArticleById = []struct {
	name            string
	requestID       string
	expectedTitle   string
	expectedContent string
	testNoID        bool
}{
	{
		name:            "get article with id 1",
		requestID:       "1",
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
		testNoID:        false,
	},
	{
		name:            "get article with id 2",
		requestID:       "2",
		expectedTitle:   testArticles[1].Title,
		expectedContent: testArticles[1].Content,
		testNoID:        false,
	},
	{
		name:            "get article with non existing id",
		requestID:       "99",
		expectedTitle:   "",
		expectedContent: "",
		testNoID:        false,
	},
	{
		name:            "get article with no id",
		requestID:       "",
		expectedTitle:   "",
		expectedContent: "",
	},
}

var testsForGetCommentsByArticleId = []struct {
	name            string
	articleId       string
	pagination      testPagination
	expectedLength  int
	expectedContent string
	expectedUserID  int
}{
	{
		name:            "default comments list",
		articleId:       "1",
		pagination:      testPagination{limit: "", page: "", order: ""},
		expectedLength:  2,
		expectedContent: testComments[0].Content,
		expectedUserID:  testComments[0].UserID,
	},
	{
		name:            "comments list with limit of 1",
		articleId:       "1",
		pagination:      testPagination{limit: "1", page: "", order: ""},
		expectedLength:  1,
		expectedContent: testComments[0].Content,
		expectedUserID:  testComments[0].UserID,
	},
	{
		name:            "comments list with page of 5",
		articleId:       "1",
		pagination:      testPagination{limit: "", page: "5", order: ""},
		expectedLength:  0,
		expectedContent: "",
	},
	{
		name:            "comments list with order of DESC",
		articleId:       "2",
		pagination:      testPagination{limit: "", page: "", order: "DESC"},
		expectedLength:  2,
		expectedContent: testComments[3].Content,
		expectedUserID:  testComments[3].UserID,
	},
	{
		name:            "comments list with mixed pagination",
		articleId:       "2",
		pagination:      testPagination{limit: "1", page: "2", order: "DESC"},
		expectedLength:  1,
		expectedContent: testComments[2].Content,
		expectedUserID:  testComments[2].UserID,
	},
	{
		name:            "comments list with incorrect pagination variables",
		articleId:       "2",
		pagination:      testPagination{limit: "test", page: "-50", order: "TestOrder"},
		expectedLength:  2,
		expectedContent: testComments[2].Content,
		expectedUserID:  testComments[2].UserID,
	},
	{
		name:            "comments list with non existing article id",
		articleId:       "99",
		pagination:      testPagination{limit: "", page: "", order: ""},
		expectedLength:  0,
		expectedContent: "",
	},
	{
		name:            "comments list with incorrect article id",
		articleId:       "",
		pagination:      testPagination{limit: "", page: "", order: ""},
		expectedLength:  0,
		expectedContent: "",
	},
}

var testsForVerifyApiKeyExists = []struct {
	name             string
	testApiKey       string
	expectedResponse bool
}{
	{
		name:             "verify with existing api key",
		testApiKey:       testUsers[0].ApiKey,
		expectedResponse: true,
	},
	{
		name:             "verify with non existing api key",
		testApiKey:       "asdf-4444-test-123",
		expectedResponse: false,
	},
	{
		name:             "verify with empty api key",
		testApiKey:       "",
		expectedResponse: false,
	},
}

var testsForInsertComment = []struct {
	name          string
	articleID     int
	userID        int
	content       string
	expectedError bool
}{
	{
		name:          "insert comment with correct data",
		articleID:     1,
		userID:        1,
		content:       "some new content 1",
		expectedError: false,
	},
	{
		name:          "insert comment with incorrect article id",
		articleID:     99,
		userID:        1,
		content:       "some new content 99",
		expectedError: true,
	},
	{
		name:          "insert comment with incorrect user id",
		articleID:     1,
		userID:        99,
		content:       "some new content user 99",
		expectedError: true,
	},
	{
		name:          "insert comment with empty content",
		articleID:     2,
		userID:        2,
		content:       "",
		expectedError: false,
	},
	{
		name:          "insert comment with empty article id",
		userID:        2,
		content:       "empty article id",
		expectedError: true,
	},
	{
		name:          "insert comment with empty user id",
		articleID:     1,
		content:       "empty user id",
		expectedError: true,
	},
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
	suite.testRepo.App = &a
	db, err := driver.ConnectSQL(testDsn, a)
	if err != nil {
		suite.Fail("could not connect to the test sql")
	}
	suite.testRepo.DB = db.SQL
	err = suite.testRepo.DB.AutoMigrate(
		&models.UserGroup{},
		&models.User{},
		&models.Article{},
		&models.Comment{},
	)
	if err != nil {
		suite.Fail("could generate test tables")
	}
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
	if err != nil {
		suite.Fail("could not drop test tables")
	}
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(databaseRequestTestSuite))
}

func (suite *databaseRequestTestSuite) TestGetArticlesList() {
	for _, t := range testsForGetArticlesList {
		req, err := generateNewGETRequest("/v1/articles", t.pagination)
		suite.Nil(err, "failed to create http request")
		list, err := suite.testRepo.GetArticlesList(req)
		if err != nil {
			suite.Fail(fmt.Sprintf("failed to load articles in test name %s", t.name))
		}
		l := len(list)
		suite.Equal(
			t.expectedLength,
			l,
			fmt.Sprintf("expected length %d is not equal to actual length %d in test name %s", t.expectedLength, l, t.name),
		)
		if t.expectedLength == 0 && l == t.expectedLength {
			continue
		}
		suite.Equal(
			t.expectedTitle,
			list[0].Title,
			fmt.Sprintf("expected and actual title are not equal in test name %s", t.name),
		)
		suite.Equal(
			t.expectedContent,
			list[0].Content,
			fmt.Sprintf("expected and actual content are not equal in test name %s", t.name),
		)
	}
}

func (suite *databaseRequestTestSuite) TestGetArticleById() {
	pagination := testPagination{
		limit: "",
		page:  "",
		order: "",
	}
	for _, t := range testsForGetArticleById {
		req, err := generateNewGETRequest("/v1/article/:id", pagination)
		suite.Nil(err, "failed to create http request")
		ctx := req.Context()
		ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
			{"id", t.requestID},
		})
		req = req.WithContext(ctx)
		a, err := suite.testRepo.GetArticleById(req)
		if err != nil {
			if t.requestID == "" {
				continue
			}
			suite.Fail(fmt.Sprintf("error when loading article by id %s, in test name %s", t.requestID, t.name))
		}
		suite.Equal(
			t.expectedTitle,
			a.Title,
			fmt.Sprintf("expected and actual title are not equal in test name %s", t.name),
		)
		suite.Equal(
			t.expectedContent,
			a.Content,
			fmt.Sprintf("expected and actual content are not equal in test name %s", t.name),
		)
	}
}

func (suite *databaseRequestTestSuite) TestGetCommentsByArticleId() {
	for _, t := range testsForGetCommentsByArticleId {
		req, err := generateNewGETRequest("/v1/article/:id/comments", t.pagination)
		suite.Nil(err, "failed to create http request")
		ctx := req.Context()
		ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
			{"id", t.articleId},
		})
		req = req.WithContext(ctx)
		list, err := suite.testRepo.GetCommentsByArticleId(req)
		if err != nil {
			if t.articleId == "" {
				continue
			}
			suite.Fail(fmt.Sprintf("failed to load articles in test name %s", t.name))
		}
		l := len(list)
		suite.Equal(
			t.expectedLength,
			l,
			fmt.Sprintf("expected length %d is not equal to actual length %d in test name %s", t.expectedLength, l, t.name),
		)
		if t.expectedLength == 0 && l == t.expectedLength {
			continue
		}
		suite.Equal(
			t.expectedUserID,
			list[0].UserID,
			fmt.Sprintf("expected and actual user id are not equal in test name %s", t.name),
		)
		suite.Equal(
			t.expectedContent,
			list[0].Content,
			fmt.Sprintf("expected and actual content are not equal in test name %s", t.name),
		)
	}
}

func (suite *databaseRequestTestSuite) TestVerifyApiKeyExists() {
	for _, t := range testsForVerifyApiKeyExists {
		v := suite.testRepo.VerifyApiKeyExists(t.testApiKey)
		suite.Equal(t.expectedResponse, v, fmt.Sprintf("expected and actual responses are not equal in test name %s", t.name))
	}
}

func (suite *databaseRequestTestSuite) TestInsertComment() {
	for _, t := range testsForInsertComment {
		var tc models.Comment
		tc.ArticleID = t.articleID
		tc.UserID = t.userID
		tc.Content = t.content
		tc.CreatedAt = time.Now()
		tc.UpdatedAt = time.Now()
		err := suite.testRepo.InsertComment(tc)
		if err != nil {
			if t.expectedError {
				continue
			}
			suite.Fail(fmt.Sprintf("error not expected in test name %s", t.name))
		}
	}
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
