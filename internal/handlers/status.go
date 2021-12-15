package handlers

import (
	"bitbucket.org/julius_liaudanskis/go-blog/models"
	"encoding/json"
	"log"
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
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		log.Println(err)
		return
	}
}
