package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// StatusHandler is a handler for app status
func (repo *Repository) StatusHandler(c *gin.Context) {
	currentStatus := models.AppStatus{
		Status:      "Available",
		Environment: repo.App.Environment,
		Version:     repo.App.AppVersion,
	}
	c.JSON(
		http.StatusAccepted,
		GetDataWrap(currentStatus),
	)
}
