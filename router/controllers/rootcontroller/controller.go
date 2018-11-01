package rootcontroller

import (
	"flag"
	"net/http"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/L-oris/yabb/repository/bucketrepository"
	"github.com/L-oris/yabb/router/httperror"
	"github.com/gorilla/mux"
)

const maxUploadSize = 2 * 1024 * 1024 // MB

type Config struct {
	Tpl tpl.Template
	Serve
	Bucket *bucketrepository.Repository
}

type Controller struct {
	Router *mux.Router
	serve  Serve
	tpl    tpl.Template
	bucket *bucketrepository.Repository
}

// New creates a new controller and registers the routes
func New(config *Config) Controller {
	c := Controller{
		serve:  config.Serve,
		tpl:    config.Tpl,
		bucket: config.Bucket,
	}

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(c.static())
	router.HandleFunc("/", c.home).Methods("GET")
	router.HandleFunc("/ping", c.ping).Methods("GET")
	router.HandleFunc("/bucket/{id}", c.serveBucketFileByID).Methods("GET")

	c.Router = router
	return c
}

// static serves static files
func (c Controller) static() http.Handler {
	var dir string
	flag.StringVar(&dir, "dir", "public/", "the directory to serve files from /public")
	flag.Parse()

	return http.StripPrefix("/static/", http.FileServer(http.Dir(dir)))
}

// home serves the Home page
func (c Controller) home(w http.ResponseWriter, req *http.Request) {
	c.tpl.Render(w, "home.gohtml", nil)
}

// ping is used for health check
func (c Controller) ping(w http.ResponseWriter, req *http.Request) {
	logger.Log.Debug("ping pong request")
	w.Write([]byte("pong"))
}

// serveBucketFileByID serves files from GC bucket
func (c Controller) serveBucketFileByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	imageID := vars["id"]

	file, err := c.bucket.Read(imageID)
	if err != nil {
		httperror.NotFound(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", http.DetectContentType(file))
	w.Write(file)
}
