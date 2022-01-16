package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
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
		expectedCode:    http.StatusInternalServerError,
		expectedTitle:   "",
		expectedContent: "",
	},
	{
		name:            "get article with no id",
		expectedCode:    http.StatusInternalServerError,
		requestID:       "",
		expectedTitle:   "",
		expectedContent: "",
	},
}

func (suite *handlersTestSuite) TestGetArticlesList() {
	for _, t := range testsForGetArticlesList {
		m := make(map[string][]models.Article)
		req, err := generateNewGETRequest("/v1/articles", t.pagination)
		suite.Nil(err, "failed to create new request")
		rr := httptest.NewRecorder()
		suite.testHandlerRepo.GetArticlesList(rr, req)
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
		err = json.Unmarshal(rr.Body.Bytes(), &m)
		suite.Nil(err, "could not unmarshal the response body into a map")
		suite.Equal(t.expectedTitle, m["data"][0].Title, fmt.Sprintf("status code should be %s but got %s", t.expectedTitle, m["data"][0].Title))
	}
}

func (suite *handlersTestSuite) TestGetArticleById() {
	var p testPagination
	for _, t := range testsForGetArticleById {
		m := make(map[string]models.Article)
		req, err := generateNewGETRequest("/v1/article/:id", p)
		suite.Nil(err, "failed to create http request")
		ctx := req.Context()
		ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
			{"id", t.requestID},
		})
		rr := httptest.NewRecorder()
		req = req.WithContext(ctx)
		suite.testHandlerRepo.GetArticleById(rr, req)
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
		err = json.Unmarshal(rr.Body.Bytes(), &m)
		suite.Nil(err, "could not unmarshal the response body into a map")
		suite.Equal(t.expectedTitle, m["data"].Title, fmt.Sprintf("status code should be %s but got %s", t.expectedTitle, m["data"].Title))
	}
}
