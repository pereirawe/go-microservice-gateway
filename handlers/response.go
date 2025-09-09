package handlers

import (
	"encoding/json"
	"net/http"
)

// ResponseHandler is a function that standardizes HTTP responses.
// It receives an http.ResponseWriter, the content, the status code, and the content type.
func ResponseHandler(w http.ResponseWriter, content []byte, statusCode int, contentType string) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	if contentType == "" {
		contentType = "application/json"
	}

	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(statusCode)

	if content != nil {
		w.Write(content)
	}
}

// CreateJSONResponse creates a JSON response with the given content and status code.
// It sets the Content-Type header to "application/json".
func CreateJSONResponse(w http.ResponseWriter, content interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if content != nil {
		json.NewEncoder(w).Encode(content)
	}
}
