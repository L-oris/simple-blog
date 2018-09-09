package tpl

import (
	"log"
	"net/http"
)

// TPL implements Template interface
type TPL struct{}

// Render gets the template, fills it with data and sends it to ResponseWriter
func (*TPL) Render(w http.ResponseWriter, name string, data interface{}) {
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
