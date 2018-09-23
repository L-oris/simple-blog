package router

import (
	"net/http"
	"os"

	"github.com/L-oris/yabb/controller/postcontroller"
	"github.com/L-oris/yabb/controller/rootcontroller"
	"github.com/L-oris/yabb/httperror"
	"github.com/L-oris/yabb/inject"
	"github.com/L-oris/yabb/inject/types"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func Mount() http.Handler {
	router := mux.NewRouter()

	router.PathPrefix("/post").Handler(
		negroni.New(negroni.Wrap(
			inject.Container.Get(types.PostController.String()).(postcontroller.Controller).Router),
		))

	router.PathPrefix("/").Handler(
		negroni.New(negroni.Wrap(
			inject.Container.Get(types.RootController.String()).(rootcontroller.Controller).Router),
		))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		httperror.NotFound(w, "Route Not Found")
	})

	return handlers.LoggingHandler(os.Stdout, router)
}
