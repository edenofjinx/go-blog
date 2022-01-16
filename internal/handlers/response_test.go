package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

var testsForResponseJson = []struct {
	name         string
	statusCode   int
	wrap         string
	data         interface{}
	expectedWrap string
	expectedCode int
}{
	{
		name:         "default response test",
		statusCode:   http.StatusOK,
		wrap:         "success",
		data:         "",
		expectedWrap: "success",
		expectedCode: http.StatusOK,
	},
	{
		name:         "response test with correct data",
		statusCode:   http.StatusOK,
		wrap:         "success",
		data:         "test data",
		expectedWrap: "success",
		expectedCode: http.StatusOK,
	},
	{
		name:         "response test with incorrect data",
		statusCode:   http.StatusInternalServerError,
		wrap:         "success",
		data:         make(chan int),
		expectedWrap: "error",
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
		log.Println(status)
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
	}
}
