package utils

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, code int, respCode, message string) {
	responseHTTPError := HTTPError{
		Code:    respCode,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(responseHTTPError)

}
