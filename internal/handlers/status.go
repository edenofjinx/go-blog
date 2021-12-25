package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
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

	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		repo.App.ErrorLog.Println(err)
	}
	w.Header().Set(AppContentType, AppJson)
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write(js)
	if err != nil {
		repo.App.ErrorLog.Println(err)
		return
	}
}
