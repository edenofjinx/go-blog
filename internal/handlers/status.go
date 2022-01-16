package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// StatusHandler is a handler for app status
func (repo *Repository) StatusHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	currentStatus := models.AppStatus{
		Status:      "Available",
		Environment: repo.App.Environment,
		Version:     repo.App.AppVersion,
	}
	err := repo.ResponseJson(w, http.StatusAccepted, currentStatus, "success")
	if err != nil {
		repo.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
}
