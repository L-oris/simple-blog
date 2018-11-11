package injector

// 'core.go' does not build any dependency with wire, therefore its functions will be excluded from generated code in 'wire_gen.go'

import (
	"database/sql"
	"net/http"

	"github.com/L-oris/yabb/foreign/template"
	"github.com/L-oris/yabb/repositories/db"
)

func provideFileServer() (func(w http.ResponseWriter, r *http.Request, name string), error) {
	return http.ServeFile, nil
}

func provideDB() *sql.DB {
	return db.BlogDB
}

func provideRenderer() (template.Renderer, error) {
	return template.Template{}, nil
}
