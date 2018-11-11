package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/L-oris/yabb/router/controllers/postcontroller"
	"github.com/L-oris/yabb/router/controllers/rootcontroller"
	"github.com/L-oris/yabb/router/httperror"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Config struct {
	PostController postcontroller.Controller
	RootController rootcontroller.Controller
}

// New creates and returns a new mux.Router, with all handled attached to it
func New(config Config) http.Handler {
	router := mux.NewRouter()

	attachHandler(router, "/post", config.PostController.Router)
	attachHandler(router, "/", config.RootController.Router)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		httperror.NotFound(w, "Route Not Found")
	})

	return handlers.LoggingHandler(os.Stdout, router)
}

func attachHandler(r *mux.Router, path string, handler http.Handler) {
	r.PathPrefix(path).Handler(
		http.StripPrefix(
			strings.TrimSuffix(path, "/"),
			handler,
		),
	)
}
