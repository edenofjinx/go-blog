package handlers

import (
	"encoding/json"
	"net/http"
)

// JsonResponse hold json response message
type JsonResponse struct {
	Message string `json:"message"`
}

// ResponseJson handles json responses
func (repo *Repository) ResponseJson(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})
	wrapper[wrap] = data
	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	w.Header().Set(AppContentType, AppJson)
	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}
	return nil
}
