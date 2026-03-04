package utils

import (
	"GoWebUser/models/dto"
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, code int, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.AuthResponse{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
