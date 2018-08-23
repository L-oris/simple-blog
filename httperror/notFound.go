package httperror

import (
	"encoding/json"
	"log"
	"net/http"
)

// NotFound writes a default 404 Not Found to the client
func NotFound(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	em := ErrorMessage{
		Status:      404,
		Message:     "Not Found",
		Description: errorMessage,
	}

	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(em)
	if err != nil {
		log.Fatalln("httperror.NotFound > encoding error:", err)
	}
}
