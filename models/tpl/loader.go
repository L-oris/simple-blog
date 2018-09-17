package tpl

import (
	"html/template"
	"path/filepath"

	"github.com/L-oris/yabb/logger"
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
		logger.Log.Fatal("get layoutFiles error: ", err)
	}

	includeFiles, err := filepath.Glob(config.TemplateIncludePath + "*.gohtml")
	if err != nil {
		logger.Log.Fatal("get includeFiles error: ", err)
	}

	mainTemplate := template.New("main")
	if mainTemplate, err = mainTemplate.Parse(mainTmpl); err != nil {
		logger.Log.Fatal("parse error: ", err)
	}

	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			logger.Log.Fatal("clone error: ", err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

	logger.Log.Debug("templates loading successful")
	bufpool = bpool.NewBufferPool(64)
	logger.Log.Debug("template buffer allocation successful")
}
