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
	apiKey := models.GenerateToken(payload.Email)
	var user models.User
	user.Password = pw
	user.Email = payload.Email
	user.Name = payload.Name
	user.ApiKey = apiKey
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
	var user models.User
	user.Email = payload.Email
	user.Name = payload.Name
	user.UpdatedAt = time.Now()
	user.ID = payload.ID
	err = repo.DB.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not save the user. Try again later."))
		return
	}
	c.JSON(http.StatusAccepted, GetSuccessMessageWrap("User has been saved."))
}

// LoginUser is a handler for user login verification
func (repo *Repository) LoginUser(c *gin.Context) {
	var payload models.UserLoginPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("An error occurred while trying to login. Try again later."))
		return
	}
	if payload.Email == "" || payload.Password == "" {
		c.JSON(http.StatusBadRequest, GetErrorMessageWrap("Missing login credentials. Try again."))
		return
	}
	user, err := repo.DB.GetUserByEmail(payload.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, GetErrorMessageWrap("User was not found."))
		return
	}
	var resp models.UserLoginResponse
	if !models.CheckPasswordHash(payload.Password, user.Password) {
		c.JSON(http.StatusForbidden, GetErrorMessageWrap("Password is incorrect. Try again."))
		return
	}
	resp.ApiKey = user.ApiKey
	log.Println(user)
	log.Println(resp)
	c.JSON(http.StatusAccepted, GetDataWrap(resp))
	return
}

// DeleteUser is a handler for deleting a user
func (repo *Repository) DeleteUser(c *gin.Context) {

}

// UpdateUserPassword is a handler for updating a users password
func (repo *Repository) UpdateUserPassword(c *gin.Context) {
	var payload models.UserPasswordPayload
	err := json.NewDecoder(c.Request.Body).Decode(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could not update the user password. Try again later."))
		return
	}
	var user models.User
	pw, err := models.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could update the user password. Try again later."))
		return
	}
	user.Password = pw
	user.UpdatedAt = time.Now()
	user.ID = payload.ID
	err = repo.DB.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GetErrorMessageWrap("Could update the user password. Try again later."))
		return
	}
	c.JSON(http.StatusAccepted, GetSuccessMessageWrap("The password has been updated."))
}
