package controller

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/L-oris/yabb/httperror"
	"github.com/L-oris/yabb/models/post"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
)

type BlogController struct {
	store map[string]post.Post
}

func NewBlogController() *BlogController {
	return &BlogController{
		store: make(map[string]post.Post),
	}
}

// Home serves the Home page
func (c BlogController) Home(w http.ResponseWriter, req *http.Request) {
	tpl.RenderTemplate(w, "home.gohtml", nil)
}

// GetAll gets all post.Post from the store
func (c BlogController) GetAll(w http.ResponseWriter, req *http.Request) {
	tpl.RenderTemplate(w, "all.gohtml", c.store)
}

// New renders the template for adding new post.Post
func (c BlogController) New(w http.ResponseWriter, req *http.Request) {
	tpl.RenderTemplate(w, "new.gohtml", c.store)
}

// Add adds a post.Post to the store
func (c BlogController) Add(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	partialPost := post.Post{
		Title:   req.Form["title"][0],
		Content: req.Form["content"][0],
	}
	newPost, err := post.GenerateFromPartial(partialPost)
	if err != nil {
		httperror.BadRequest(w, "Invalid data provided")
		return
	}

	c.store[newPost.ID] = newPost
	tpl.RenderTemplate(w, "byID.gohtml", newPost)
}

// GetByID gets a post.Post by ID from store
func (c BlogController) GetByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fPostID := vars["id"]
	fPost := c.store[fPostID]
	if fPost == post.EmptyPost {
		errorMessage := "Post " + fPostID + " not found"
		httperror.NotFound(w, errorMessage)
		return
	}

	tpl.RenderTemplate(w, "byID.gohtml", fPost)
}

// UpdateByID accepts a partial post.Post and update its fields
func (c BlogController) UpdateByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postID := vars["id"]
	currentPost := c.store[postID]
	if currentPost == post.EmptyPost {
		errMessage := "Post " + postID + " not found"
		httperror.NotFound(w, errMessage)
		return
	}

	bsJSON, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln("controller.UpdateById > reading error:", err)
	}
	defer req.Body.Close()

	newPartialPost, err := post.FromJSON(bsJSON)
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

// DeleteByID deletes a post.Post by ID
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

// LoggingMiddleware logs all incoming requests
func (c BlogController) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("controller.LoggingMiddleware:", req.Method, req.RequestURI)
		next.ServeHTTP(w, req)
	})
}
