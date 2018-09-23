package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/L-oris/yabb/inject/types"
	"github.com/L-oris/yabb/router/controllers/postcontroller"
	"github.com/L-oris/yabb/router/controllers/rootcontroller"
	"github.com/L-oris/yabb/router/httperror"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sarulabs/di"
)

// Mount creates and returns a new mux.Router, with all handled attached to it
func Mount(ctn di.Container) http.Handler {
	router := mux.NewRouter()

	attachHandler(router, "/post", ctn.Get(types.PostController.String()).(postcontroller.Controller).Router)
	attachHandler(router, "/", ctn.Get(types.RootController.String()).(rootcontroller.Controller).Router)
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
