package handlers

import (
	"github.com/gin-gonic/gin"
)

// JsonResponse hold json response message
type JsonResponse struct {
	Message string `json:"message"`
}

func GetSuccessMessageWrap(data string) gin.H {
	return gin.H{
		"success": JsonResponse{
			data,
		},
	}
}

func GetErrorMessageWrap(data string) gin.H {
	return gin.H{
		"error": JsonResponse{
			data,
		},
	}
}

func GetDataWrap(data interface{}) gin.H {
	return gin.H{
		"data": data,
	}
}
