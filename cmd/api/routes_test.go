package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (suite *TestMainPackage) TestProtectedRoutes() {
	router := httprouter.New()
	secure := alice.New(verifyApiKey)
	protectedRoutes(router, &secure)
	h, _, _ := router.Lookup(http.MethodGet, "/v1/article/1")
	suite.NotNil(h, "/v1/article/:id route should exist")
}

func (suite *TestMainPackage) TestUnprotectedRoutes() {
	router := httprouter.New()
	unprotectedRoutes(router)
	h, _, _ := router.Lookup(http.MethodGet, "/v1/status")
	suite.NotNil(h, "/v1/status route should exist")
}

func (suite *TestMainPackage) TestRoutes() {
	mux := routes()
	switch v := mux.(type) {
	case http.Handler:
		//correct handler, do nothing
	default:
		suite.Fail(fmt.Sprintf("type is not http.Handler but is %T", v))
	}
}
