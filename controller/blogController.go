package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/L-oris/mongoRestAPI/httperror"
	"github.com/L-oris/mongoRestAPI/models"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
)

type BlogController struct {
	store map[string]models.Post
}

func NewBlogController() *BlogController {
	return &BlogController{
		store: make(map[string]models.Post),
	}
}

// Home serves the Home page
func (c BlogController) Home(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello world"))
}

// GetAll gets all models.Post from the store
func (c BlogController) GetAll(w http.ResponseWriter, req *http.Request) {
	if len(c.store) == 0 {
		w.Write([]byte("The store is empty"))
		return
	}

	err := json.NewEncoder(w).Encode(c.store)
	if err != nil {
		log.Fatalln("controller.GetAll > encoding error:", err)
	}
}

// Add adds a models.Post to the store
func (c BlogController) Add(w http.ResponseWriter, req *http.Request) {
	bsJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("controller.Add > reading error:", err)
	}
	defer req.Body.Close()

	newPost, err := models.PostFromJSON(bsJSON)
	if err != nil {
		httperror.BadRequest(w, "Invalid JSON")
		return
	}

	if !models.IsValidPost(newPost) {
		httperror.BadRequest(w, "Bad Post")
		return
	}

	c.store[newPost.ID] = newPost
	w.Write([]byte("OK"))
}

// GetByID gets a models.Post by ID from store
func (c BlogController) GetByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postID := vars["id"]
	post := c.store[postID]
	if post == models.EmptyPost {
		errorMessage := "Post " + postID + " not found"
		httperror.NotFound(w, errorMessage)
		return
	}

	err := json.NewEncoder(w).Encode(post)
	if err != nil {
		log.Fatalln("controller.GetAll > encoding error:", err)
	}
}

// UpdateByID accepts a partial models.Post and update its fields
func (c BlogController) UpdateByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postID := vars["id"]
	currentPost := c.store[postID]
	if currentPost == models.EmptyPost {
		errMessage := "Post " + postID + " not found"
		httperror.NotFound(w, errMessage)
		return
	}

	bsJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("controller.UpdateById > reading error:", err)
	}
	defer req.Body.Close()

	newPartialPost, err := models.PostFromJSON(bsJSON)
	if err != nil {
		httperror.BadRequest(w, "Invalid JSON")
		return
	}

	if err = mergo.Merge(&currentPost, newPartialPost, mergo.WithOverride); err != nil {
		log.Println("controller.UpdateById > invalid post provided:", err)
		errMessage := "Invalid post provided"
		httperror.BadRequest(w, errMessage)
		return
	}

	c.store[currentPost.ID] = currentPost
	w.Write([]byte("OK"))
}

// DeleteByID deletes a models.Post by ID
// It doesn't care if the Post exists or not
func (c BlogController) DeleteByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postID := vars["id"]
	delete(c.store, postID)
	w.Write([]byte("OK"))
}

// RouteNotFound handles requests to routes not implemented
func (c BlogController) RouteNotFound(w http.ResponseWriter, req *http.Request) {
	httperror.NotFound(w, "Route Not Found")
}
