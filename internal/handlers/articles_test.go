package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
