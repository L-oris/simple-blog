package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/L-oris/mongoRestAPI/httperror"
	"github.com/L-oris/mongoRestAPI/models"
	"github.com/julienschmidt/httprouter"
)

type BlogController struct {
	store map[string]models.Post
}

func NewBlogController() *BlogController {
	return &BlogController{
		store: make(map[string]models.Post),
	}
}

func (c BlogController) GetAll(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if len(c.store) == 0 {
		w.Write([]byte("The store is empty"))
		return
	}

	jsonStore, err := json.Marshal(c.store)
	if err != nil {
		log.Fatalln("models.GetAll > marshaling error:", err)
	}
	w.Write(jsonStore)
}

func (c BlogController) Add(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	bsJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("models.Add > reading error:", err)
	}
	defer req.Body.Close()

	var newPost models.Post
	err = json.Unmarshal(bsJSON, &newPost)
	if err != nil {
		log.Println("models.Add > unmarshal error:", err)
		httperror.BadRequest(w, "Invalid JSON")
		return
	}

	if !models.IsValidPost(newPost) {
		log.Println("models.Add > bad post received")
		httperror.BadRequest(w, "Bad Post")
		return
	}

	c.store[newPost.ID] = newPost
	w.Write([]byte("OK"))
}
