package helper

import (
	"encoding/json"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(payload)
}
