package router

import (
	"net/http"
	"os"

	"github.com/L-oris/yabb/inject/types"
	"github.com/L-oris/yabb/router/controllers/postcontroller"
	"github.com/L-oris/yabb/router/controllers/rootcontroller"
	"github.com/L-oris/yabb/router/httperror"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sarulabs/di"
)

func Mount(ctn di.Container) http.Handler {
	router := mux.NewRouter()

	router.PathPrefix("/post").Handler(
		ctn.Get(types.PostController.String()).(postcontroller.Controller).Router,
	)

	router.PathPrefix("/").Handler(
		ctn.Get(types.RootController.String()).(rootcontroller.Controller).Router,
	)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		httperror.NotFound(w, "Route Not Found")
	})

	return handlers.LoggingHandler(os.Stdout, router)
}
