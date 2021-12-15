package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"net/http"
)

func (repo *Repository) StatusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := models.AppStatus{
		Status:      "Available",
		Environment: repo.App.Environment,
		Version:     repo.App.AppVersion,
	}

	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		repo.App.ErrorLog.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write(js)
	if err != nil {
		repo.App.ErrorLog.Println(err)
		return
	}
}
