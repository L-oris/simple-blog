package httperror

import (
	"encoding/json"
	"log"
	"net/http"
)

// InternalServer writes a default 500 Internal Server to the client
func InternalServer(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	if errorMessage == "" {
		errorMessage = "Internal Server Error"
	}
	em := ErrorMessage{
		Status:      http.StatusBadRequest,
		Message:     "Internal Server Error",
		Description: errorMessage,
	}

	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(em)
	if err != nil {
		log.Fatalln("httperror.InternalServer > encoding error:", err)
	}
}
