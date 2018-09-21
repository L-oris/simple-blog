package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/L-oris/yabb/controller/postcontroller"
	"github.com/L-oris/yabb/controller/rootcontroller"
	"github.com/L-oris/yabb/httperror"
	"github.com/L-oris/yabb/inject"
	"github.com/L-oris/yabb/models/env"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/L-oris/yabb/repository/postrepository"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/post").Handler(negroni.New(
		negroni.Wrap(postcontroller.New(&postcontroller.Config{
			PathPrefix: "/post",
			Repository: inject.Container.Get("postrepository").(*postrepository.Repository),
			Tpl:        inject.Container.Get("templates").(*tpl.TPL),
		}).Router)))

	router.PathPrefix("/").Handler(negroni.New(negroni.Wrap(rootcontroller.New(
		&rootcontroller.Config{
			PathPrefix: "/",
			Tpl:        inject.Container.Get("templates").(*tpl.TPL),
			Serve:      inject.Container.Get("fileserver").(func(w http.ResponseWriter, r *http.Request, fileName string)),
		}).Router)))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		httperror.NotFound(w, "Route Not Found")
	})

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	server := &http.Server{
		Addr:         ":" + env.Vars.Port,
		Handler:      loggedRouter,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Fatal(server.ListenAndServe())
}
