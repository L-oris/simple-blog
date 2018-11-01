package tpl

import (
	"net/http"

	"github.com/L-oris/yabb/logger"
)

// TPL implements Template interface
type TPL struct{}

// Render gets the template, fills it with data and sends it to ResponseWriter
func (tpl *TPL) Render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		logger.Log.Errorf("template %s does not exist, rendering default error template", name)
		tpl.Render(w, "error.gohtml", nil)
		return
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	if err := tmpl.Execute(buf, data); err != nil {
		logger.Log.Errorf("cannot execute template %s, rendering default error template", name)
		tpl.Render(w, "error.gohtml", nil)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
