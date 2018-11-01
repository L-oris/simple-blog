package httperror

import (
	"encoding/json"
	"net/http"

	"github.com/L-oris/yabb/logger"
)

// InternalServer writes a 500 Internal Server to the client
func InternalServer(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	em := ErrorMessage{
		Status:      http.StatusBadRequest,
		Message:     "Internal Server Error",
		Description: errorMessage,
	}

	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(em)
	if err != nil {
		logger.Log.Error("encoding error:", err)
		w.Write([]byte("Internal Server Error"))
	}
}
