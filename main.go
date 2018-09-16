package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/L-oris/yabb/models/db"
	"github.com/L-oris/yabb/repository/postrepository"

	"github.com/L-oris/yabb/controller/postcontroller"
	"github.com/L-oris/yabb/controller/rootcontroller"
	"github.com/L-oris/yabb/httperror"
	"github.com/L-oris/yabb/models/env"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	router := mux.NewRouter()

	router.PathPrefix("/post").Handler(negroni.New(
		negroni.Wrap(postcontroller.New(&postcontroller.Config{
			PathPrefix: "/post",
			Repository: postrepository.New(),
			Tpl:        &tpl.TPL{},
		}).Router)))

	router.PathPrefix("/").Handler(negroni.New(negroni.Wrap(rootcontroller.New(
		&rootcontroller.Config{
			DB:         db.BlogDB,
			PathPrefix: "/",
			Tpl:        &tpl.TPL{},
			ServeFile:  http.ServeFile,
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
