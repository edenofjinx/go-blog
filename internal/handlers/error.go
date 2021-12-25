package handlers

import (
	"net/http"
)

// ErrorHandler handles error responses
func (repo *Repository) ErrorHandler(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	jsError := JsonResponse{
		Message: err.Error(),
	}
	writeError := repo.ResponseJson(w, statusCode, jsError, "error")
	if writeError != nil {
		repo.App.ErrorLog.Println(writeError)
		return
	}
}
