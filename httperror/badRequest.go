package httperror

import (
	"encoding/json"
	"log"
	"net/http"
)

// BadRequest writes a default 400 Bad Request to the client
func BadRequest(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	em := ErrorMessage{
		Status:      http.StatusBadRequest,
		Message:     "Bad Request",
		Description: errorMessage,
	}

	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(em)
	if err != nil {
		log.Fatalln("httperror.BadRequest > encoding error:", err)
	}
}
