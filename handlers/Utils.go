package handlers

import (
	"encoding/json"
	"github.com/kcapp/odds-api/models"
	"net/http"
)

// SetHeaders will set the default headers used by all requests
func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func JSONError(w http.ResponseWriter, err string, code int) {
	e := models.Error{
		Message: err,
		Status:  code,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(e)
}

func JSONSuccess(w http.ResponseWriter, err string, code int) {
	e := models.Error{
		Message: err,
		Status:  code,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(e)
}
