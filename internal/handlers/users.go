package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// CreateUser is a handler for creating a user
func (repo *Repository) CreateUser(c *gin.Context) {
	var payload models.UserPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the user. Try again later."))
		return
	}
	pw, err := models.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the user. Try again later."))
		return
	}
	var user models.User
	user.Password = pw
	user.Email = payload.Email
	user.Name = payload.Name
	user.UpdatedAt = time.Now()
	user.CreatedAt = time.Now()
	//TODO add default registered user group
	user.GroupID = 2
	err = repo.DB.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the user. Try again later."))
		return
	}
	c.JSON(http.StatusAccepted, GetSuccessMessageWrap("User has been saved."))
}

// UpdateUser is a handler for updating a user
func (repo *Repository) UpdateUser(c *gin.Context) {
	var payload models.UserUpdatePayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the user. Try again later."))
		return
	}
	log.Println(payload)
	var user models.User
	user.Email = payload.Email
	user.Name = payload.Name
	user.UpdatedAt = time.Now()
	user.ID = payload.ID
	log.Println(user)
	err = repo.DB.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the user. Try again later."))
		return
	}
	c.JSON(http.StatusAccepted, GetSuccessMessageWrap("User has been saved."))
}

// LoginUser is a handler for user login verification
func (repo *Repository) LoginUser(c *gin.Context) {

}

// DeleteUser is a handler for deleting a user
func (repo *Repository) DeleteUser(c *gin.Context) {

}
