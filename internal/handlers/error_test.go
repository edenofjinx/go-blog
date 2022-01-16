package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
)

var testsForErrorHandler = []struct {
	name            string
	statusCode      int
	statusMessage   string
	expectedCode    int
	expectedMessage string
}{
	{
		name:            "default error handler test",
		statusCode:      http.StatusBadRequest,
		statusMessage:   "Testing. Bad request",
		expectedCode:    http.StatusBadRequest,
		expectedMessage: "Testing. Bad request",
	},
	{
		name:            "error handler test without statusCode",
		statusCode:      0,
		statusMessage:   "Testing. No statusCode",
		expectedCode:    http.StatusBadRequest,
		expectedMessage: "Testing. No statusCode",
	},
}

func (suite *handlersTestSuite) TestErrorHandler() {
	for _, t := range testsForErrorHandler {
		m := make(map[string]JsonResponse)
		rr := httptest.NewRecorder()
		if t.statusCode != 0 {
			suite.testHandlerRepo.ErrorHandler(rr, errors.New(t.statusMessage), t.statusCode)
		} else {
			suite.testHandlerRepo.ErrorHandler(rr, errors.New(t.statusMessage))
		}
		status := rr.Code
		suite.Equal(t.expectedCode, status, fmt.Sprintf("status code should be %d but got %d", t.expectedCode, status))
		err := json.Unmarshal(rr.Body.Bytes(), &m)
		suite.Nil(err, "could not unmarshal the response body into a struct")
		suite.Equal(t.expectedMessage, m["error"].Message, "expected message is not equal to actual")
	}
}
