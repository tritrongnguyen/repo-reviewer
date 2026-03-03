package helper

import (
	"encoding/json"
	"net/http"
)

type APIResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// SuccessResponse returns a successful API response.
func SuccessResponse[T any](code int, data T, message string) APIResponse[T] {
	return APIResponse[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse returns a failure API response.
func ErrorResponse(code int, err string) APIResponse[any] {
	return APIResponse[any]{
		Code:    code,
		Message: err,
	}
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(payload)
}
