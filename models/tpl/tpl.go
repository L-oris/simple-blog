package tpl

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/oxtoacart/bpool"
)

type TemplateConfig struct {
	TemplateLayoutPath  string
	TemplateIncludePath string
}

func init() {
	loadConfiguration()
	loadTemplates()
}

// TPL holds a reference to all templates
var TPL *template.Template

var templateConfig TemplateConfig
var templates map[string]*template.Template
var bufpool *bpool.BufferPool
var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`

func loadConfiguration() {
	templateConfig.TemplateLayoutPath = "templates/layouts/"
	templateConfig.TemplateIncludePath = "templates/"
}

func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles, err := filepath.Glob(templateConfig.TemplateLayoutPath + "*.gohtml")
	if err != nil {
		log.Fatalln("tpl.loadTemplates > get layoutFiles error:", err)
	}

	includeFiles, err := filepath.Glob(templateConfig.TemplateIncludePath + "*.gohtml")
	if err != nil {
		log.Fatalln("tpl.loadTemplates > get includeFiles error:", err)
	}

	mainTemplate := template.New("main")
	if mainTemplate, err = mainTemplate.Parse(mainTmpl); err != nil {
		log.Fatalln("tpl.loadTemplates > error:", err)
	}

	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatalln("tpl.loadTemplates > error:", err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

	log.Println("tpl.loadTemplates > templates loading successful")
	bufpool = bpool.NewBufferPool(64)
	log.Println("tpl.loadTemplates > buffer allocation successful")
}

// RenderTemplate gets the template, fills it with data and sends it to ResponseWriter
func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		log.Fatalf("tpl.RenderTemplate > template %s does not exist", name)
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	if err := tmpl.Execute(buf, data); err != nil {
		log.Fatalln("tpl.RenderTemplate > cannot execute template", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
