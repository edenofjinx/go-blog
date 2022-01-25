package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

var testsForGetArticlesList = []struct {
	name          string
	expectedCode  int
	expectedTitle string
	pagination    testPagination
}{
	{
		name:          "default articles list",
		expectedCode:  http.StatusAccepted,
		expectedTitle: testArticles[0].Title,
		pagination: testPagination{
			limit: "",
			page:  "",
			order: "",
		},
	},
	{
		name:          "articles list with pagination",
		expectedCode:  http.StatusAccepted,
		expectedTitle: testArticles[0].Title,
		pagination: testPagination{
			limit: "1",
			page:  "2",
			order: "DESC",
		},
	},
	{
		name:          "articles list with incorrect pagination",
		expectedCode:  http.StatusAccepted,
		expectedTitle: testArticles[0].Title,
		pagination: testPagination{
			limit: "9999",
			page:  "-50",
			order: "test",
		},
	},
}

var testsForGetArticleById = []struct {
	name            string
	requestID       string
	expectedCode    int
	expectedTitle   string
	expectedContent string
}{
	{
		name:          "default article by id",
		requestID:     "1",
		expectedCode:  http.StatusAccepted,
		expectedTitle: testArticles[0].Title,
	},
	{
		name:            "get article with non existing id",
		requestID:       "99",
		expectedCode:    http.StatusAccepted,
		expectedTitle:   "",
		expectedContent: "",
	},
	{
		name:            "get article with non existing id string",
		requestID:       "test",
		expectedCode:    http.StatusBadRequest,
		expectedTitle:   "",
		expectedContent: "",
	},
	{
		name:            "get article with no id",
		expectedCode:    http.StatusBadRequest,
		requestID:       "",
		expectedTitle:   "",
		expectedContent: "",
	},
	{
		name:            "get article with no string as id",
		expectedCode:    http.StatusBadRequest,
		requestID:       "test",
		expectedTitle:   "",
		expectedContent: "",
	},
}

var testsForSaveArticle = []struct {
	name         string
	jsonData     string
	expectedCode int
}{
	{
		name:         "article with incorrect json",
		jsonData:     "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 1,\"id\": 1,,,}",
		expectedCode: http.StatusInternalServerError,
	},
	{
		name:         "article with incorrect base64 image png",
		jsonData:     "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBOR w0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 1,\"id\": 1}",
		expectedCode: http.StatusInternalServerError,
	},
	{
		name:         "article update with incorrect user id",
		jsonData:     "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 99,\"id\": 1}",
		expectedCode: http.StatusInternalServerError,
	},
	{
		name:         "new article with correct json",
		jsonData:     "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 1}",
		expectedCode: http.StatusAccepted,
	},
	{
		name:         "new article with incorrect json",
		jsonData:     "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 999}",
		expectedCode: http.StatusInternalServerError,
	},
	{
		name:         "article update with correct json",
		jsonData:     "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 1,\"id\": 1}",
		expectedCode: http.StatusAccepted,
	},
	{
		name:         "article update with incorrect json",
		jsonData:     "{\"content\": \"<p>test</p><img src='data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAFCAYAAACNbyblAAAAHElEQVQI12P4//8/w38GIAXDIBKE0DHxgljNBAAO9TXL0Y4OHwAAAABJRU5ErkJggg==' />\",\"title\": \"test title\",\"user_id\": 999,\"id\": 1}",
		expectedCode: http.StatusInternalServerError,
	},
}

var testsForDeleteArticle = []struct {
	name         string
	articleId    string
	expectedCode int
}{
	{
		name:         "article id as set to string",
		articleId:    "test",
		expectedCode: http.StatusBadRequest,
	},
	{
		name:         "article with given id is not available",
		articleId:    "9584",
		expectedCode: http.StatusAccepted,
	},
	{
		name:         "successful article deletion",
		articleId:    "1",
		expectedCode: http.StatusAccepted,
	},
}

func (suite *handlersTestSuite) TestGetArticlesList() {
	for _, t := range testsForGetArticlesList {
		m := make(map[string][]models.Article)
		c, rr := generateNewGETRequest("/v1/articles/", t.pagination)
		suite.testHandlerRepo.GetArticlesList(c)
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
		err := json.Unmarshal(rr.Body.Bytes(), &m)
		suite.Nil(err, "could not unmarshal the response body into a map")
		suite.Equal(t.expectedTitle, m["data"][0].Title, fmt.Sprintf("status code should be %s but got %s", t.expectedTitle, m["data"][0].Title))
	}
}

func (suite *handlersTestSuite) TestGetArticleById() {
	var p testPagination
	for _, t := range testsForGetArticleById {
		m := make(map[string]models.Article)
		c, rr := generateNewGETRequest("/v1/article/:id", p)
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: t.requestID,
			},
		}
		suite.testHandlerRepo.GetArticleById(c)
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
		if t.expectedCode == status && t.expectedContent == "" {
			continue
		}
		err := json.Unmarshal(rr.Body.Bytes(), &m)
		suite.Nil(err, "could not unmarshal the response body into a map")
		suite.Equal(t.expectedTitle, m["data"].Title, fmt.Sprintf("status code should be %s but got %s", t.expectedTitle, m["data"].Title))
	}
}

func (suite *handlersTestSuite) TestSaveArticle() {
	for _, t := range testsForSaveArticle {
		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)
		req := &http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Path: "/v1/article/save",
			},
			Body: io.NopCloser(strings.NewReader(t.jsonData)),
		}
		c.Request = req
		suite.testHandlerRepo.SaveArticle(c)
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d: %s", t.expectedCode, status, t.name))
	}
}

func (suite *handlersTestSuite) TestDeleteArticle() {
	var p testPagination
	for _, t := range testsForDeleteArticle {
		c, rr := generateNewGETRequest("/v1/comment/delete/:id", p)
		c.Params = []gin.Param{
			{
				Key:   "id",
				Value: t.articleId,
			},
		}
		suite.testHandlerRepo.DeleteArticle(c)
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
	}
}
