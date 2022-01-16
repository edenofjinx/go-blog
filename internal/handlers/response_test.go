package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

var testsForResponseJson = []struct {
	name         string
	statusCode   int
	wrap         string
	data         interface{}
	expectedCode int
}{
	{
		name:         "default response test",
		statusCode:   http.StatusOK,
		wrap:         "success",
		data:         "",
		expectedCode: http.StatusOK,
	},
	{
		name:         "response test with correct data",
		statusCode:   http.StatusOK,
		wrap:         "success",
		data:         "test data",
		expectedCode: http.StatusOK,
	},
	{
		name:         "response test with incorrect data",
		statusCode:   http.StatusInternalServerError,
		wrap:         "success",
		data:         make(chan int),
		expectedCode: http.StatusInternalServerError,
	},
}

func (suite *handlersTestSuite) TestResponseJson() {
	for _, t := range testsForResponseJson {
		rr := httptest.NewRecorder()
		err := suite.testHandlerRepo.ResponseJson(rr, t.statusCode, t.data, t.wrap)
		if err != nil {
			suite.testHandlerRepo.ErrorHandler(rr, err, t.statusCode)
		}
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
	}
}
