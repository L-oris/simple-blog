package rootcontroller

import (
	"database/sql"
	"flag"
	"net/http"

	"github.com/L-oris/yabb/foreign/template"
	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/repositories/bucketrepository"
	"github.com/L-oris/yabb/router/httperror"
	"github.com/gorilla/mux"
)

const maxUploadSize = 2 * 1024 * 1024 // MB

type Config struct {
	Renderer template.Renderer
	Serve    func(w http.ResponseWriter, r *http.Request, fileName string)
	Bucket   *bucketrepository.Repository
	DB       *sql.DB
}

type Controller struct {
	Router   *mux.Router
	serve    func(w http.ResponseWriter, r *http.Request, fileName string)
	renderer template.Renderer
	bucket   *bucketrepository.Repository
	db       *sql.DB
}

// New creates a new controller and registers the routes
func New(config *Config) Controller {
	c := Controller{
		serve:    config.Serve,
		renderer: config.Renderer,
		bucket:   config.Bucket,
		db:       config.DB,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", c.home).Methods("GET")
	router.HandleFunc("/ping", c.ping).Methods("GET")
	router.HandleFunc("/pingDB", c.pingDB).Methods("GET")
	router.PathPrefix("/static/").Handler(c.static())
	router.HandleFunc("/bucket/{id}", c.serveBucketFileByID).Methods("GET")

	c.Router = router
	return c
}

func NewWire(config Config) Controller {
	c := Controller{
		serve:    config.Serve,
		renderer: config.Renderer,
		bucket:   config.Bucket,
		db:       config.DB,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", c.home).Methods("GET")
	router.HandleFunc("/ping", c.ping).Methods("GET")
	router.HandleFunc("/pingDB", c.pingDB).Methods("GET")
	router.PathPrefix("/static/").Handler(c.static())
	router.HandleFunc("/bucket/{id}", c.serveBucketFileByID).Methods("GET")

	c.Router = router
	return c
}

// ping is used for health check on the server
func (c Controller) ping(w http.ResponseWriter, req *http.Request) {
	logger.Log.Debug("ping pong request - server")
	w.Write([]byte("pong - server"))
}

// ping is used for health check on the database
func (c Controller) pingDB(w http.ResponseWriter, req *http.Request) {
	logger.Log.Debug("ping pong request - DB")
	if err := c.db.Ping(); err != nil {
		logger.Log.Errorf("db ping error: ", err)
		httperror.InternalServer(w, "failed ping connection to db")
		return
	}
	w.Write([]byte("pong - DB"))
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
	c.renderer.Render(w, "home.gohtml", nil)
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
