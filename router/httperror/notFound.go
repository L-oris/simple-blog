package httperror

import (
	"encoding/json"
	"net/http"

	"github.com/L-oris/yabb/logger"
)

// NotFound writes a 404 Not Found to the client
func NotFound(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	em := ErrorMessage{
		Status:      http.StatusNotFound,
		Message:     "Not Found",
		Description: errorMessage,
	}

	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(em)
	if err != nil {
		logger.Log.Error("encoding error:", err)
		w.Write([]byte("Not Found"))
	}
}
