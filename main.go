package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", Welcome)
	// router.GET("/posts", GetAll)
	// router.POST("/post", Add)
	// router.GET("/post/:id", GetById)
	// router.PUT("/post/:id", UpdateById)
	// router.DELETE("/post/:id", DeleteById)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Welcome(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Write([]byte("Hello world"))
}
