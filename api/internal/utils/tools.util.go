package utils

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("No se pudo escribir la respuesta JSON: %v", err)
	}
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "El cuerpo JSON no es válido."})
		return false
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		JsonResponse(w, http.StatusBadRequest, map[string]string{"error": "El cuerpo debe contener un único objeto JSON."})
		return false
	}
	return true
}

func EncodeJSON(w http.ResponseWriter, value any) {
	JsonResponse(w, http.StatusOK, value)
}

type Validator interface {
	GetErrors() map[string]string
}

func ValidateAndRespond(w http.ResponseWriter, isValid bool, validator Validator) bool {
	if !isValid {
		JsonResponse(w, http.StatusBadRequest, map[string]any{
			"error":   "Hay errores de validación.",
			"details": validator.GetErrors(),
		})
		return false
	}
	return true
}
