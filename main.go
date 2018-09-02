package main

import (
	"log"
	"net/http"
	"time"

	"github.com/L-oris/yabb/controller"
	"github.com/gorilla/mux"
)

func main() {
	blogController := controller.NewBlogController()
	router := mux.NewRouter()
	router.Use(blogController.LoggingMiddleware)
	router.HandleFunc("/", blogController.Home).Methods("GET")
	router.HandleFunc("/posts", blogController.GetAll).Methods("GET")
	router.HandleFunc("/add", blogController.New).Methods("GET")
	router.HandleFunc("/post", blogController.Add).Methods("POST")
	router.HandleFunc("/post/{id}", blogController.GetByID).Methods("GET")
	router.HandleFunc("/post/{id}/edit", blogController.EditByID).Methods("GET")
	router.HandleFunc("/post/{id}/edit", blogController.UpdateByID).Methods("POST")
	router.HandleFunc("/post/{id}", blogController.DeleteByID).Methods("DELETE")

	router.NotFoundHandler = http.HandlerFunc(blogController.RouteNotFound)
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Fatal(server.ListenAndServe())
}

// TODO: handle .favicon request
