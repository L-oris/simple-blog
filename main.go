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
	router.GET("/", Welcome)
	router.GET("/posts", blogController.GetAll)
	router.POST("/post", blogController.Add)
	router.GET("/post/:id", blogController.GetByID)
	router.PUT("/post/:id", blogController.UpdateByID)
	// router.DELETE("/post/:id", DeleteById)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Welcome(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Write([]byte("Hello world"))
}
