package models

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// TPL holds a reference to all templates
var TPL *template.Template

func init() {
	TPL = template.Must(findAndParseTemplates("templates"))
}

func findAndParseTemplates(rootDir string) (*template.Template, error) {
	funcMap := template.FuncMap{}
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".gohtml") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			t, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}
