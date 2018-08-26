package main

import (
	"log"
	"net/http"
	"time"

	"github.com/L-oris/mongoRestAPI/controller"
	"github.com/gorilla/mux"
)

func main() {
	blogController := controller.NewBlogController()
	router := mux.NewRouter()
	router.HandleFunc("/", blogController.Home).Methods("GET")
	router.HandleFunc("/posts", blogController.GetAll).Methods("GET")
	router.HandleFunc("/post", blogController.Add).Methods("POST")
	router.HandleFunc("/post/{id}", blogController.GetByID).Methods("GET")
	router.HandleFunc("/post/{id}", blogController.UpdateByID).Methods("PUT")
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
