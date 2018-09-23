package postcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/models/post"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/L-oris/yabb/repository/postrepository"
	"github.com/L-oris/yabb/router/httperror"
	"github.com/gorilla/mux"
)

type Config struct {
	PathPrefix string
	Repository *postrepository.Repository
	Tpl        tpl.Template
}

type postControllerStore map[string]post.Post

type Controller struct {
	repository *postrepository.Repository
	tpl        tpl.Template
	Router     *mux.Router
}

// New creates a new controller and registers the routes
func New(config *Config) Controller {
	c := Controller{
		repository: config.Repository,
		tpl:        config.Tpl,
	}

	router := mux.NewRouter()
	routes := router.PathPrefix(config.PathPrefix).Subrouter()
	routes.HandleFunc("/ping", c.ping).Methods("GET")
	routes.HandleFunc("/all", c.renderAll).Methods("GET")
	routes.HandleFunc("/new", c.renderNew).Methods("GET")
	routes.HandleFunc("/{id}", c.renderByID).Methods("GET")
	routes.HandleFunc("/{id}/update", c.renderUpdateByID).Methods("GET")
	routes.HandleFunc("/new", c.new).Methods("POST")
	routes.HandleFunc("/{id}/update", c.updateByID).Methods("POST")
	routes.HandleFunc("/{id}/delete", c.deleteByID).Methods("POST")

	c.Router = router
	return c
}

// ping checks db connection
func (c Controller) ping(w http.ResponseWriter, req *http.Request) {
	if err := c.repository.Ping(); err != nil {
		httperror.InternalServer(w, "db not connected")
		return
	}
	w.Write([]byte("pong"))
}
