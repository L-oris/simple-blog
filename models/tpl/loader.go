package tpl

import (
	"html/template"
	"log"
	"path/filepath"

	"github.com/oxtoacart/bpool"
)

func init() {
	loadConfiguration()
	loadTemplates()
}

type templateConfig struct {
	TemplateLayoutPath  string
	TemplateIncludePath string
}

var config templateConfig
var templates map[string]*template.Template
var bufpool *bpool.BufferPool
var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`

func loadConfiguration() {
	config.TemplateLayoutPath = "templates/layouts/"
	config.TemplateIncludePath = "templates/"
}

func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles, err := filepath.Glob(config.TemplateLayoutPath + "*.gohtml")
	if err != nil {
		log.Fatalln("tpl.loadTemplates > get layoutFiles error:", err)
	}

	includeFiles, err := filepath.Glob(config.TemplateIncludePath + "*.gohtml")
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
