package rootcontroller

import (
	"flag"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/L-oris/yabb/repository/bucketrepository"
	"github.com/L-oris/yabb/router/httperror"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
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
	router.HandleFunc("/upload", c.uploadGet).Methods("GET")
	router.HandleFunc("/upload", c.uploadPost).Methods("POST")
	router.HandleFunc("/bucket/{id}", c.serveBucket).Methods("GET")

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

func (c Controller) uploadGet(w http.ResponseWriter, req *http.Request) {
	c.tpl.Render(w, "upload.gohtml", nil)
}

func (c Controller) uploadPost(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, maxUploadSize)
	if err := req.ParseMultipartForm(maxUploadSize); err != nil {
		logger.Log.Debug("uploaded file is too big: %s", err.Error())
		httperror.BadRequest(w, "file provided is too large")
		return
	}

	multipartFile, _, err := req.FormFile("postImage")
	if err != nil {
		logger.Log.Error("could not get form from template: %s", err.Error())
		httperror.InternalServer(w, "invalid template form")
		return
	}
	defer multipartFile.Close()

	fileBytes, err := ioutil.ReadAll(multipartFile)
	if err != nil {
		logger.Log.Debug("invalid file uploaded: %s", err.Error())
		httperror.BadRequest(w, "invalid file provided")
		return
	}

	contentType := http.DetectContentType(fileBytes)
	if ok := checkContentType(contentType); !ok {
		httperror.BadRequest(w, "invalid fileType provided")
		return
	}

	fileEndings, _ := mime.ExtensionsByType(contentType)
	fileName := uuid.NewV4().String()
	logger.Log.Debug("ContentType: %s, File: %s", contentType, fileName+fileEndings[0])

	err = c.bucket.Write(fileName, fileBytes)
	if err != nil {
		httperror.InternalServer(w, "cannot save file")
		return
	}
	w.Write([]byte("uploading ok"))
}

func checkContentType(fileType string) bool {
	if fileType != "image/jpeg" &&
		fileType != "image/jpg" &&
		fileType != "image/gif" &&
		fileType != "image/png" {
		return false
	}
	return true
}

func (c Controller) serveBucket(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	imageID := vars["id"]

	file, err := c.bucket.Read(imageID)
	if err != nil {
		httperror.NotFound(w, fmt.Sprintf("file %s not found", imageID))
	}

	w.Header().Set("Content-Type", http.DetectContentType(file))
	w.Write(file)
}
