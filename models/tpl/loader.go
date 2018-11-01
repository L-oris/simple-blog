package tpl

import (
	"errors"
	"html/template"
	"path/filepath"

	"github.com/L-oris/yabb/logger"
	"github.com/oxtoacart/bpool"
)

func init() {
	loadConfiguration()
	if err := loadTemplates(); err != nil {
		logger.Log.Fatal(err.Error())
	}
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

func loadTemplates() error {
	defaultError := errors.New("could not load templates")

	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles, err := filepath.Glob(config.TemplateLayoutPath + "*.gohtml")
	if err != nil {
		logger.Log.Critical("get layoutFiles error: ", err)
		return defaultError
	}

	includeFiles, err := filepath.Glob(config.TemplateIncludePath + "*.gohtml")
	if err != nil {
		logger.Log.Critical("get includeFiles error: ", err)
		return defaultError
	}

	mainTemplate := template.New("main")
	if mainTemplate, err = mainTemplate.Parse(mainTmpl); err != nil {
		logger.Log.Critical("parse error: ", err)
		return defaultError
	}

	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			logger.Log.Critical("clone error: ", err)
			return defaultError
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

	logger.Log.Debug("templates loading successful")
	bufpool = bpool.NewBufferPool(64)
	logger.Log.Debug("template buffer allocation successful")
	return nil
}
