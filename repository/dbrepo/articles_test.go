package dbrepo

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"fmt"
)

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
	requestID       int
	expectedTitle   string
	expectedContent string
	testNoID        bool
}{
	{
		name:            "get article with id 1",
		requestID:       1,
		expectedTitle:   testArticles[0].Title,
		expectedContent: testArticles[0].Content,
		testNoID:        false,
	},
	{
		name:            "get article with id 2",
		requestID:       2,
		expectedTitle:   testArticles[1].Title,
		expectedContent: testArticles[1].Content,
		testNoID:        false,
	},
	{
		name:            "get article with non existing id",
		requestID:       99,
		expectedTitle:   "",
		expectedContent: "",
		testNoID:        false,
	},
	{
		name:            "get article with no id",
		requestID:       -500,
		expectedTitle:   "",
		expectedContent: "",
	},
}

var testsForSaveArticle = []struct {
	name          string
	jsonData      string
	expectedError bool
}{
	{
		name:          "save article",
		jsonData:      "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 1}",
		expectedError: false,
	},
	{
		name:          "save article with empty data",
		jsonData:      "{\"content\": \"\",\"title\": \"\",\"user_id\": 0}",
		expectedError: true,
	},
	{
		name:          "save article with incorrect data",
		jsonData:      "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 999}",
		expectedError: true,
	},
}

var testsForUpdateArticle = []struct {
	name          string
	jsonData      string
	expectedError bool
}{
	{
		name:          "update article",
		jsonData:      "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 1, \"id\": 1}",
		expectedError: false,
	},
	{
		name:          "save article with empty data",
		jsonData:      "{\"content\": \"\",\"title\": \"\"}",
		expectedError: true,
	},
	{
		name:          "save article with incorrect data",
		jsonData:      "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 999, \"id\": 1}",
		expectedError: true,
	},
}

var testsForDeleteArticle = []struct {
	name          string
	articleId     int
	expectedError bool
}{
	{
		name:          "delete article with incorrect id",
		articleId:     9999,
		expectedError: false,
	},
	{
		name:          "delete article with 0 id",
		articleId:     0,
		expectedError: false,
	},
	{
		name:          "delete article",
		articleId:     1,
		expectedError: false,
	},
}

func (suite *databaseRequestTestSuite) TestGetArticlesList() {
	for _, t := range testsForGetArticlesList {
		c := generateNewGETRequest("/v1/articles", t.pagination)
		list, err := suite.testRepo.GetArticlesList(c)
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
	for _, t := range testsForGetArticleById {
		ar, err := suite.testRepo.GetArticleById(t.requestID)
		if err != nil {
			if t.requestID == 0 {
				continue
			}
			suite.Fail(fmt.Sprintf("error when loading article by id %d, in test name %s", t.requestID, t.name))
		}
		suite.Equal(
			t.expectedTitle,
			ar.Title,
			fmt.Sprintf("expected and actual title are not equal in test name %s", t.name),
		)
		suite.Equal(
			t.expectedContent,
			ar.Content,
			fmt.Sprintf("expected and actual content are not equal in test name %s", t.name),
		)
	}
}

func (suite *databaseRequestTestSuite) TestSaveArticle() {
	for _, t := range testsForSaveArticle {
		var payload models.ArticlePayload
		var article models.Article
		err := json.Unmarshal([]byte(t.jsonData), &payload)
		suite.Nil(err, "should not be nil")
		article.UserID = payload.UserID
		article.Title = payload.Title
		article.Content = payload.Content
		err = suite.testRepo.SaveArticle(article)
		if t.expectedError {
			suite.Error(err, "should be an error")
		} else {
			suite.Nil(err, "should not be a nil")
		}
	}
}

func (suite *databaseRequestTestSuite) TestUpdateArticle() {
	for _, t := range testsForUpdateArticle {
		var payload models.ArticlePayload
		var article models.Article
		err := json.Unmarshal([]byte(t.jsonData), &payload)
		suite.Nil(err, "should not be nil")
		article.UserID = payload.UserID
		article.Title = payload.Title
		article.Content = payload.Content
		article.ID = payload.ID
		err = suite.testRepo.UpdateArticle(article)
		if t.expectedError {
			suite.Error(err, "should be an error")
		} else {
			suite.Nil(err, "should not be a nil")
		}
	}
}

func (suite *databaseRequestTestSuite) TestDeleteArticle() {
	for _, t := range testsForDeleteArticle {
		err := suite.testRepo.DeleteArticle(t.articleId)
		if t.expectedError {
			suite.Error(err, "should be an error")
		} else {
			suite.Nil(err, "should not be a nil")
		}
	}
}
