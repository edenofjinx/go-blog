package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/internal/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

// verifyApiKey middleware to verify if api key provided in the header exists
func verifyApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get(handlers.AppApiKey)
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, handlers.GetErrorMessageWrap("Api key was not provided"))
			return
		}
		exists := handlers.Repo.DB.VerifyApiKeyExists(apiKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, handlers.GetErrorMessageWrap("Incorrect API key provided"))
			return
		}
	}
}
