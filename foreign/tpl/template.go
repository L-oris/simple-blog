package tpl

import (
	"net/http"
)

// Template allows templates to be rendered by name on 'http.ResponseWriter'
type Template interface {
	Render(w http.ResponseWriter, name string, data interface{})
}
