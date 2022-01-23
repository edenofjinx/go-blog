package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func (suite *handlersTestSuite) TestStatusHandler() {
	m := make(map[string]models.AppStatus)
	c, rr := generateNewGETRequest("/v1/status/", testPagination{})
	suite.testHandlerRepo.StatusHandler(c)
	status := rr.Code
	suite.Equal(http.StatusAccepted, status, fmt.Sprintf("status code should be %d but got %d", http.StatusAccepted, status))
	err := json.Unmarshal(rr.Body.Bytes(), &m)
	suite.Nil(err, "could not unmarshal the response body into a struct")
	suite.Equal("test", m["data"].Environment, "expected environment is not equal to actual")
	suite.Equal("test", m["data"].Version, "expected version is not equal to actual")
}
