package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetArticlesList handler to get articles list
func (repo *Repository) GetArticlesList(c *gin.Context) {
	articles, err := repo.DB.GetArticlesList(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not get article list"))
		return
	}
	c.JSON(http.StatusAccepted, GetDataWrap(articles))
}

// GetArticleById handler to get article data by article id
func (repo *Repository) GetArticleById(c *gin.Context) {
	articles, err := repo.DB.GetArticleById(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not get article by id."))
		return
	}
	c.JSON(http.StatusAccepted, GetDataWrap(articles))
}
