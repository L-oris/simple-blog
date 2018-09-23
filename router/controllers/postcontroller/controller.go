package postcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/models/tpl"
	"github.com/L-oris/yabb/repository/postrepository"
	"github.com/L-oris/yabb/router/httperror"
	"github.com/gorilla/mux"
)

type Config struct {
	Repository *postrepository.Repository
	Tpl        tpl.Template
}

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
	router.HandleFunc("/ping", c.ping).Methods("GET")
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

// ping checks db connection
func (c Controller) ping(w http.ResponseWriter, req *http.Request) {
	if err := c.repository.Ping(); err != nil {
		httperror.InternalServer(w, "db not connected")
		return
	}
	w.Write([]byte("pong"))
}
