package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/L-oris/yabb/controller"
	"github.com/L-oris/yabb/httperror"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/posts").Handler(negroni.New(
		negroni.Wrap(controller.NewBlogController("/posts")),
	))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		httperror.NotFound(w, "Route Not Found")
	})
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      loggedRouter,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Fatal(server.ListenAndServe())
}
