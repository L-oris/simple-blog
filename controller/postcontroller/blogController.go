package postcontroller

import (
	"log"
	"net/http"

	"github.com/L-oris/yabb/httperror"
	"github.com/L-oris/yabb/models/post"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
)

type BlogController struct {
	store  map[string]post.Post
	router *mux.Router
}

func New(pathPrefix string) *mux.Router {
	c := BlogController{
		store:  make(map[string]post.Post),
		router: mux.NewRouter(),
	}

	routes := c.router.PathPrefix(pathPrefix).Subrouter()
	routes.HandleFunc("/", c.home).Methods("GET")
	routes.HandleFunc("/all", c.getAll).Methods("GET")
	routes.HandleFunc("/add", c.new).Methods("GET")
	routes.HandleFunc("/post", c.add).Methods("POST")
	routes.HandleFunc("/post/{id}", c.getByID).Methods("GET")
	routes.HandleFunc("/post/{id}/edit", c.editByID).Methods("GET")
	routes.HandleFunc("/post/{id}/edit", c.updateByID).Methods("POST")
	routes.HandleFunc("/post/{id}", c.deleteByID).Methods("DELETE")
	routes.HandleFunc("/favicon.ico", c.favicon).Methods("GET")

	return c.router
}

// Home serves the Home page
func (c BlogController) home(w http.ResponseWriter, req *http.Request) {
	tpl.RenderTemplate(w, "home.gohtml", nil)
}

// GetAll gets all Post from the store
func (c BlogController) getAll(w http.ResponseWriter, req *http.Request) {
	tpl.RenderTemplate(w, "all.gohtml", c.store)
}

// New renders the template for adding new Post
func (c BlogController) new(w http.ResponseWriter, req *http.Request) {
	tpl.RenderTemplate(w, "new.gohtml", nil)
}

// Add adds a Post to the store
func (c BlogController) add(w http.ResponseWriter, req *http.Request) {
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

// GetByID gets a Post by ID from store
func (c BlogController) getByID(w http.ResponseWriter, req *http.Request) {
	post, ok := c.getPostByIDFromStore(w, req)
	if !ok {
		return
	}
	tpl.RenderTemplate(w, "byID.gohtml", post)
}

func (c BlogController) editByID(w http.ResponseWriter, req *http.Request) {
	post, ok := c.getPostByIDFromStore(w, req)
	if !ok {
		return
	}

	tpl.RenderTemplate(w, "edit.gohtml", post)
}

// UpdateByID accepts a partial Post and update its fields
func (c BlogController) updateByID(w http.ResponseWriter, req *http.Request) {
	storePost, ok := c.getPostByIDFromStore(w, req)
	if !ok {
		return
	}

	req.ParseForm()
	newPartialPost := post.Post{
		Title:   req.Form["title"][0],
		Content: req.Form["content"][0],
	}

	if err := mergo.Merge(&storePost, newPartialPost, mergo.WithOverride); err != nil {
		log.Println("controller.UpdateById > invalid post provided:", err)
		errMessage := "Invalid post provided"
		httperror.BadRequest(w, errMessage)
		return
	}

	c.store[storePost.ID] = storePost
	tpl.RenderTemplate(w, "byID.gohtml", storePost)
}

// DeleteByID deletes a Post by ID
// It doesn't care if the Post exists or not
func (c BlogController) deleteByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postID := vars["id"]
	delete(c.store, postID)
	w.Write([]byte("OK"))
}

func (c BlogController) favicon(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "public/favicon.ico")
}
