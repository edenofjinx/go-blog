package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func (suite *TestMainPackage) TestSetProtectedRoutes() {
	mux := gin.Default()
	protected := mux.Group("/")
	setProtectedRoutes(protected)
	r := mux.Routes()
	found := false
	for _, p := range r {
		if p.Path == "/v1/article/:id" {
			found = true
			break
		}
	}
	suite.Equal(true, found, "/v1/article/:id endpoint should exist")
}

func (suite *TestMainPackage) TestSetUnprotectedRoutes() {
	mux := gin.Default()
	unprotected := mux.Group("/")
	setUnprotectedRoutes(unprotected)
	r := mux.Routes()
	found := false
	for _, p := range r {
		if p.Path == "/v1/status" {
			found = true
			break
		}
	}
	suite.Equal(true, found, "/v1/status endpoint should exist")
}

func (suite *TestMainPackage) TestRoutes() {
	mux := routes()
	r := mux.Routes()
	l := len(r)
	suite.NotEqual(0, l, "mux should have routes assigned")
}

func (suite *TestMainPackage) TestNotExistingRoute() {
	v := routes()
	req := httptest.NewRequest(http.MethodGet, "/test/test/v1", nil)
	res := httptest.NewRecorder()
	v.ServeHTTP(res, req)
	suite.Equal(http.StatusNotFound, res.Code, "response codes should be equal")
}
