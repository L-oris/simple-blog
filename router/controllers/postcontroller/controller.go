package postcontroller

import (
	"github.com/L-oris/yabb/foreign/template"
	"github.com/L-oris/yabb/repository/bucketrepository"
	"github.com/L-oris/yabb/repository/postrepository"
	"github.com/gorilla/mux"
)

type Config struct {
	Repository *postrepository.Repository
	Template   template.Template
	Bucket     *bucketrepository.Repository
}

type Controller struct {
	Router     *mux.Router
	bucket     *bucketrepository.Repository
	repository *postrepository.Repository
	template   template.Template
}

// New creates a new controller and registers the routes
func New(config *Config) Controller {
	c := Controller{
		repository: config.Repository,
		template:   config.Template,
		bucket:     config.Bucket,
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
