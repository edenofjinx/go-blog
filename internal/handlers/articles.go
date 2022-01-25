package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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
	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, GetErrorMessageWrap("Could not load article with id."))
	}
	articles, err := repo.DB.GetArticleById(articleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not get article by id."))
		return
	}
	c.JSON(http.StatusAccepted, GetDataWrap(articles))
}

// SaveArticle handler to save/update an article
func (repo *Repository) SaveArticle(c *gin.Context) {
	var payload models.ArticlePayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the article. Try again later."))
		return
	}
	content, er := repo.parseImageTags(payload.Content)
	if er != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Error saving images. Try again later."))
		return
	}
	var am models.Article
	am.Title = payload.Title
	am.Content = content
	am.UpdatedAt = time.Now()
	am.UserID = payload.UserID
	if payload.ID != 0 {
		am.ID = payload.ID
		err := repo.DB.UpdateArticle(am)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Error updating an article. Try again later."))
			return
		}
		c.JSON(http.StatusAccepted, GetSuccessMessageWrap("An article has been updated."))
		return
	}
	am.CreatedAt = time.Now()
	err = repo.DB.SaveArticle(am)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the article. Try again later."))
		return
	}
	c.JSON(http.StatusAccepted, GetSuccessMessageWrap("Article has been saved."))
}

// DeleteArticle is a handler for deleting an article
func (repo *Repository) DeleteArticle(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, GetErrorMessageWrap("Incorrect article id provided."))
		return
	}
	err = repo.DB.DeleteArticle(articleId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not delete an article."))
		return
	}
	c.JSON(http.StatusAccepted, GetSuccessMessageWrap("An article has been deleted."))
}
