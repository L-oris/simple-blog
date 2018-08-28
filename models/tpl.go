package models

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// TPL holds a reference to all templates
var TPL *template.Template

func init() {
	TPL = template.Must(parseTemplates())
}

func parseTemplates() (*template.Template, error) {
	templ := template.New("")
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".gohtml") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return err
	})

	if err != nil {
		// TODO: better error message, wrap it
		return nil, err
	}

	return templ, nil
}
