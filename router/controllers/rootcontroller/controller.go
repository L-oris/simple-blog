package rootcontroller

import (
	"context"
	"flag"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/L-oris/yabb/router/httperror"
	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/gcsblob"
	"github.com/google/go-cloud/gcp"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

const uploadPath = "./tmp"
const maxUploadSize = 2 * 1024 * 1024 // MB

type Config struct {
	Tpl tpl.Template
	Serve
}

type Controller struct {
	Router *mux.Router
	serve  Serve
	tpl    tpl.Template
}

// New creates a new controller and registers the routes
func New(config *Config) Controller {
	c := Controller{
		serve: config.Serve,
		tpl:   config.Tpl,
	}

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(c.static())
	router.HandleFunc("/", c.home).Methods("GET")
	router.HandleFunc("/ping", c.ping).Methods("GET")
	router.HandleFunc("/favicon.ico", c.favicon).Methods("GET")
	router.HandleFunc("/upload", c.uploadGet).Methods("GET")
	router.HandleFunc("/upload", c.uploadPost).Methods("POST")

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

func (c Controller) favicon(w http.ResponseWriter, req *http.Request) {
	c.serve(w, req, "public/favicon.ico")
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
	newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
	logger.Log.Debug("ContentType: %s, File: %s\n", contentType, newPath)

	// here bucket code

	bucket, err := setupGCP(CTX, "yabb")
	if err != nil {
		logger.Log.Fatalf("setup bucket error: %s", err.Error())
	}
	bucketWriter, err := bucket.NewWriter(CTX, fileName, nil)
	if err != nil {
		logger.Log.Fatalf("create bucketWriter error: %s", err.Error())
	}
	if _, err := bucketWriter.Write(fileBytes); err != nil {
		logger.Log.Fatalf("write to bucket error: %s", err.Error())
	}
	if err := bucketWriter.Close(); err != nil {
		logger.Log.Fatalf("close bucket error: %s", err.Error())
	}

	// until here bucket code

	newFile, err := os.Create(newPath)
	if err != nil {
		logger.Log.Error("could not create new empty file: %s", err.Error())
		httperror.InternalServer(w, "")
		return
	}
	defer newFile.Close()

	if _, err := newFile.Write(fileBytes); err != nil {
		logger.Log.Error("could not write bytes[] into new empty file: %s", err.Error())
		httperror.InternalServer(w, "")
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

var CTX context.Context

func init() {
	CTX = context.Background()
}

func setupGCP(ctx context.Context, bucket string) (*blob.Bucket, error) {
	credentials, err := gcp.DefaultCredentials(CTX)
	if err != nil {
		return nil, err
	}

	client, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(credentials))
	if err != nil {
		return nil, err
	}

	return gcsblob.OpenBucket(CTX, bucket, client)
}
