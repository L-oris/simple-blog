package rootcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/gorilla/mux"
)

type Config struct {
	PathPrefix string
	Tpl        tpl.Template
	ServeFile  serveFile
}

type rootController struct {
	tpl tpl.Template
	serveFile
	Router *mux.Router
}

// New creates a new controller and registers the routes
func New(config *Config) rootController {
	c := rootController{
		tpl:       config.Tpl,
		serveFile: config.ServeFile,
	}

	router := mux.NewRouter()
	routes := router.PathPrefix(config.PathPrefix).Subrouter()
	routes.HandleFunc("/", c.home).Methods("GET")
	routes.HandleFunc("/ping", c.ping).Methods("GET")
	routes.HandleFunc("/favicon.ico", c.favicon).Methods("GET")

	c.Router = router
	return c
}

// home serves the Home page
func (c rootController) home(w http.ResponseWriter, req *http.Request) {
	c.tpl.Render(w, "home.gohtml", nil)
}

// ping is used for health check
func (c rootController) ping(w http.ResponseWriter, req *http.Request) {
	logger.Log.Debug("ping pong request")
	w.Write([]byte("pong"))
}

func (c rootController) favicon(w http.ResponseWriter, req *http.Request) {
	c.serveFile(w, req, "public/favicon.ico")
}
