package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetCommentsByArticleId handler to get comments by article id
func (repo *Repository) GetCommentsByArticleId(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, GetErrorMessageWrap("Could not load article with id."))
	}
	comments, err := repo.DB.GetCommentsByArticleId(articleId, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetErrorMessageWrap("Could not get comments by article id."))
		return
	}
	c.JSON(http.StatusAccepted, GetDataWrap(comments))
}

// SaveComment saves comment into the database
func (repo *Repository) SaveComment(c *gin.Context) {
	var payload models.CommentPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the comment. Try again later."))
		return
	}
	content, er := repo.parseImageTags(payload.Content)
	if er != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Error saving images. Try again later."))
		return
	}
	var comment models.Comment
	comment.Content = content
	comment.UserID = payload.UserID
	comment.ArticleID = payload.ArticleID
	comment.UpdatedAt = time.Now()
	if payload.ID != 0 {
		comment.ID = payload.ID
		err := repo.DB.UpdateComment(comment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Error updating comment. Try again later."))
			return
		}
		c.JSON(http.StatusAccepted, GetSuccessMessageWrap("A comment has been updated."))
		return
	}
	comment.CreatedAt = time.Now()
	err = repo.DB.SaveComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the comment"))
		return
	}
	c.JSON(http.StatusOK, GetSuccessMessageWrap("Comment has been successfully saved."))
}

// DeleteComment is a handler for deleting a comment
func (repo *Repository) DeleteComment(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, GetErrorMessageWrap("Incorrect comment id provided."))
		return
	}
	err = repo.DB.DeleteComment(commentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not delete a comment."))
		return
	}
	c.JSON(http.StatusAccepted, GetSuccessMessageWrap("A comment has been deleted."))
}
