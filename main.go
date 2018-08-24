package main

import (
	"log"
	"net/http"

	"github.com/L-oris/mongoRestAPI/controller"
	"github.com/julienschmidt/httprouter"
)

func main() {
	blogController := controller.NewBlogController()
	router := httprouter.New()
	router.GET("/", blogController.Home)
	router.GET("/posts", blogController.GetAll)
	router.POST("/post", blogController.Add)
	router.GET("/post/:id", blogController.GetByID)
	router.PUT("/post/:id", blogController.UpdateByID)
	router.DELETE("/post/:id", blogController.DeleteByID)

	log.Fatal(http.ListenAndServe(":8080", router))
}
