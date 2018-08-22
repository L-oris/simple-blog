package httperror

// ErrorMessage is a default format for API messages
type ErrorMessage struct {
	Status      int    `json:"status"`
	Message     string `json:"message"`
	Description string `json:"description"`
}
