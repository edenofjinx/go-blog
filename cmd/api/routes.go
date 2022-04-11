package main

import (
	"bitbucket.org/julius_liaudanskis/go-blog/internal/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

// routes holds data of all available routes for the app
func routes() *gin.Engine {
	mux := gin.Default()
	protected := mux.Group("/")
	protected.Use(verifyApiKey())
	setProtectedRoutes(protected)
	unprotected := mux.Group("/")
	setUnprotectedRoutes(unprotected)
	mux.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, handlers.GetErrorMessageWrap("The page was not found."))
	})
	mux.Static("/static", "../../static")
	return mux
}

// setUnprotectedRoutes holds routes that are not protected by an api key
func setUnprotectedRoutes(rg *gin.RouterGroup) {
	// status
	rg.GET("/v1/status", handlers.Repo.StatusHandler)
	// users
	rg.POST("/v1/user/create", handlers.Repo.CreateUser)
	rg.POST("/v1/user/login", handlers.Repo.LoginUser)
}

// setProtectedRoutes holds routes that are protected by an api key
func setProtectedRoutes(rg *gin.RouterGroup) {
	// articles
	rg.GET("/v1/articles", handlers.Repo.GetArticlesList)
	rg.GET("/v1/article/:id", handlers.Repo.GetArticleById)
	rg.POST("/v1/article/save", handlers.Repo.SaveArticle)
	rg.DELETE("/v1/article/delete/:id", handlers.Repo.DeleteArticle)
	// comments
	rg.GET("/v1/article/:id/comments", handlers.Repo.GetCommentsByArticleId)
	rg.POST("/v1/comment/save", handlers.Repo.SaveComment)
	rg.DELETE("/v1/comment/delete/:id", handlers.Repo.DeleteComment)
	// users
	rg.PUT("/v1/user/update", handlers.Repo.UpdateUser)
	rg.PUT("/v1/user/update/password", handlers.Repo.UpdateUserPassword)
	rg.PUT("/v1/user/update/key", handlers.Repo.UpdateUserApiKey)
	rg.PUT("/v1/user/update/group", handlers.Repo.UpdateUserGroup)
	rg.DELETE("/v1/user/delete", handlers.Repo.DeleteUser)
}
