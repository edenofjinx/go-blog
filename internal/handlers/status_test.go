package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/http/httptest"
)

func (suite *handlersTestSuite) TestStatusHandler() {
	m := make(map[string]models.AppStatus)
	req, err := http.NewRequest("GET", "/v1/status", nil)
	params := httprouter.Params{}
	suite.Nil(err, "failed to create new request")
	rr := httptest.NewRecorder()
	suite.testHandlerRepo.StatusHandler(rr, req, params)
	status := rr.Code
	suite.Equal(http.StatusAccepted, status, fmt.Sprintf("status code should be %d but got %d", http.StatusAccepted, status))
	err = json.Unmarshal(rr.Body.Bytes(), &m)
	log.Println(m["success"].Status)
	suite.Nil(err, "could not unmarshal the response body into a struct")
	suite.Equal("test", m["success"].Environment, "expected environment is not equal to actual")
	suite.Equal("test", m["success"].Version, "expected version is not equal to actual")
}
