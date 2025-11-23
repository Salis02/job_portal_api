package http

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct{
	Message string `json:"message"`
	Data interface{}`json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}){
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, msg string){
	JSON(w, status, ErrorResponse{Message: msg})
}

func OK(w http.ResponseWriter, status int, msg string){
	JSON(w, status, ErrorResponse{Message: msg})
}
