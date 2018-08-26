package models

import (
	"html/template"
)

// TPL holds a reference to all templates
var TPL *template.Template

func init() {
	TPL = template.Must(template.ParseGlob("templates/*.gohtml"))
}
