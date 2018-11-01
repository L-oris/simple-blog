package httperror

import (
	"encoding/json"
	"net/http"

	"github.com/L-oris/yabb/logger"
)

// BadRequest writes a 400 Bad Request to the client
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
		logger.Log.Error("encoding error:", err)
		w.Write([]byte("Bad Request"))
	}
}
