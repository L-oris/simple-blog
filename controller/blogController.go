package controller

import (
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
	tplName := "home.gohtml"
	if err := models.TPL.ExecuteTemplate(w, tplName, nil); err != nil {
		log.Fatalln("controller.Home > cannot execute", tplName, "template:", err)
	}
}

// GetAll gets all models.Post from the store
func (c BlogController) GetAll(w http.ResponseWriter, req *http.Request) {
	if err := models.TPL.ExecuteTemplate(w, "all.gohtml", c.store); err != nil {
		log.Fatalln("controller.GetAll > executing template error:", err)
	}
}

// New renders the template for adding new models.Post
func (c BlogController) New(w http.ResponseWriter, req *http.Request) {
	if err := models.TPL.ExecuteTemplate(w, "new.gohtml", c.store); err != nil {
		log.Fatalln("controller.New > executing template error:", err)
	}
}

// Add adds a models.Post to the store
func (c BlogController) Add(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	partialPost := models.Post{
		Title:   req.Form["title"][0],
		Content: req.Form["content"][0],
	}
	newPost, err := models.GeneratePost(partialPost)
	if err != nil {
		httperror.BadRequest(w, "Invalid data provided")
		return
	}

	c.store[newPost.ID] = newPost
	if err := models.TPL.ExecuteTemplate(w, "byID.gohtml", newPost); err != nil {
		log.Fatalln("controller.Add > executing template error:", err)
	}
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

	if err := models.TPL.ExecuteTemplate(w, "byID.gohtml", post); err != nil {
		log.Fatalln("controller.GetByID > executing template error:", err)
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

	newPartialPost, err := models.FromJSON(bsJSON)
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

// LoggingMiddleware logs all incoming requests
func (c BlogController) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("controller.LoggingMiddleware:", req.Method, req.RequestURI)
		next.ServeHTTP(w, req)
	})
}
