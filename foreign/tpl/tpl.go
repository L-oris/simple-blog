package tpl

import (
	"net/http"

	"github.com/L-oris/yabb/logger"
)

// TPL implements Template interface
type TPL struct{}

// Render gets the template, fills it with data and sends it to ResponseWriter
func (tpl *TPL) Render(w http.ResponseWriter, templateName string, data interface{}) {
	tmpl, ok := templates[templateName]
	if !ok {
		logger.Log.Errorf("template %s does not exist", templateName)
		handleRenderError(w, templateName)
		return
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	if err := tmpl.Execute(buf, data); err != nil {
		logger.Log.Errorf("cannot execute template %s", templateName)
		handleRenderError(w, templateName)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func handleRenderError(w http.ResponseWriter, templateName string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("An error occurred rendering the template"))
}
