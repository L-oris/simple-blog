package template

import (
	"net/http"
)

// Renderer allows templates to be rendered by name on 'http.ResponseWriter'
type Renderer interface {
	Render(w http.ResponseWriter, templateName string, data interface{})
}
