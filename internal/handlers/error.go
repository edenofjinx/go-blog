package handlers

import (
	"encoding/json"
	"net/http"
)

type jsonError struct {
	Message string `json:"message"`
}

func (repo *Repository) ErrorHandler(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	jsError := jsonError{
		Message: err.Error(),
	}
	writeError := writeJSON(w, statusCode, jsError, "error")
	if writeError != nil {
		repo.App.ErrorLog.Println(writeError)
		return
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
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
