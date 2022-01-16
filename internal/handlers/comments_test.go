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

var testsForGetCommentsByArticleId = []struct {
	name            string
	requestID       string
	expectedCode    int
	expectedContent string
	pagination      testPagination
}{
	{
		name:            "default articles list",
		requestID:       "1",
		expectedCode:    http.StatusAccepted,
		expectedContent: testComments[0].Content,
		pagination: testPagination{
			limit: "",
			page:  "",
			order: "",
		},
	},
	{
		name:            "articles list with pagination",
		requestID:       "1",
		expectedCode:    http.StatusAccepted,
		expectedContent: testComments[0].Content,
		pagination: testPagination{
			limit: "1",
			page:  "2",
			order: "DESC",
		},
	},
	{
		name:            "comments list with incorrect pagination",
		requestID:       "2",
		expectedCode:    http.StatusAccepted,
		expectedContent: testComments[2].Content,
		pagination: testPagination{
			limit: "9999",
			page:  "-50",
			order: "test",
		},
	},
	{
		name:            "comments list with non existing id",
		requestID:       "99",
		expectedCode:    http.StatusAccepted,
		expectedContent: "",
		pagination: testPagination{
			limit: "",
			page:  "",
			order: "",
		},
	},
	{
		name:            "comments list with non empty id",
		requestID:       "",
		expectedCode:    http.StatusBadRequest,
		expectedContent: "",
		pagination: testPagination{
			limit: "",
			page:  "",
			order: "",
		},
	},
	{
		name:            "comments list with string as id",
		requestID:       "test",
		expectedCode:    http.StatusBadRequest,
		expectedContent: "",
		pagination: testPagination{
			limit: "",
			page:  "",
			order: "",
		},
	},
}

var testsForSaveComment = []struct {
	name string
}{
	{},
}

func (suite *handlersTestSuite) TestGetCommentsByArticleId() {
	for _, t := range testsForGetCommentsByArticleId {
		m := make(map[string][]models.Comment)
		req, err := generateNewGETRequest("/v1/article/:id/comments", t.pagination)
		suite.Nil(err, "failed to create http request")
		ctx := req.Context()
		ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
			{"id", t.requestID},
		})
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		suite.testHandlerRepo.GetCommentsByArticleId(rr, req)
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
		if t.expectedCode == http.StatusBadRequest && t.expectedContent == "" {
			continue
		}
		err = json.Unmarshal(rr.Body.Bytes(), &m)
		suite.Nil(err, "could not unmarshal the response body into a map")
		if len(m) == 0 && t.expectedContent != "" {
			suite.Equal(t.expectedContent, m["data"][0].Content, fmt.Sprintf("status code should be %s but got %s", t.expectedContent, m["data"][0].Content))
		}
	}
}
