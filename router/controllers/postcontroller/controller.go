package postcontroller

import (
	"github.com/L-oris/yabb/foreign/template"
	"github.com/L-oris/yabb/services/postservice"
	"github.com/gorilla/mux"
)

type Config struct {
	Renderer template.Renderer
	Service  *postservice.Service
}

type Controller struct {
	Router   *mux.Router
	renderer template.Renderer
	service  *postservice.Service
}

// New creates a new controller and registers the routes
func New(config *Config) Controller {
	c := Controller{
		renderer: config.Renderer,
		service:  config.Service,
	}

	router := mux.NewRouter()
	router.HandleFunc("/all", c.renderAll).Methods("GET")
	router.HandleFunc("/new", c.renderNew).Methods("GET")
	router.HandleFunc("/{id}", c.renderByID).Methods("GET")
	router.HandleFunc("/{id}/update", c.renderUpdateByID).Methods("GET")
	router.HandleFunc("/new", c.new).Methods("POST")
	router.HandleFunc("/{id}/update", c.updateByID).Methods("POST")
	router.HandleFunc("/{id}/delete", c.deleteByID).Methods("POST")

	c.Router = router
	return c
}

// NewWire creates a new controller and registers the routes
func NewWire(config Config) Controller {
	c := Controller{
		renderer: config.Renderer,
		service:  config.Service,
	}

	router := mux.NewRouter()
	router.HandleFunc("/all", c.renderAll).Methods("GET")
	router.HandleFunc("/new", c.renderNew).Methods("GET")
	router.HandleFunc("/{id}", c.renderByID).Methods("GET")
	router.HandleFunc("/{id}/update", c.renderUpdateByID).Methods("GET")
	router.HandleFunc("/new", c.new).Methods("POST")
	router.HandleFunc("/{id}/update", c.updateByID).Methods("POST")
	router.HandleFunc("/{id}/delete", c.deleteByID).Methods("POST")

	c.Router = router
	return c
}
