package controller

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

func NewBlogController(pathPrefix string) *mux.Router {
	c := BlogController{
		store:  make(map[string]post.Post),
		router: mux.NewRouter(),
	}

	routes := c.router.PathPrefix(pathPrefix).Subrouter()
	routes.Use(c.LoggingMiddleware)
	routes.HandleFunc("/home", c.Home)
	routes.HandleFunc("/", c.Home).Methods("GET")
	routes.HandleFunc("/all", c.GetAll).Methods("GET")
	routes.HandleFunc("/add", c.New).Methods("GET")
	routes.HandleFunc("/post", c.Add).Methods("POST")
	routes.HandleFunc("/post/{id}", c.GetByID).Methods("GET")
	routes.HandleFunc("/post/{id}/edit", c.EditByID).Methods("GET")
	routes.HandleFunc("/post/{id}/edit", c.UpdateByID).Methods("POST")
	routes.HandleFunc("/post/{id}", c.DeleteByID).Methods("DELETE")
	routes.HandleFunc("/favicon.ico", c.Favicon).Methods("GET")

	return c.router
}

// Home serves the Home page
func (c BlogController) Home(w http.ResponseWriter, req *http.Request) {
	tpl.RenderTemplate(w, "home.gohtml", nil)
}

// GetAll gets all Post from the store
func (c BlogController) GetAll(w http.ResponseWriter, req *http.Request) {
	tpl.RenderTemplate(w, "all.gohtml", c.store)
}

// New renders the template for adding new Post
func (c BlogController) New(w http.ResponseWriter, req *http.Request) {
	tpl.RenderTemplate(w, "new.gohtml", nil)
}

// Add adds a Post to the store
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

// GetByID gets a Post by ID from store
func (c BlogController) GetByID(w http.ResponseWriter, req *http.Request) {
	post, ok := c.getPostByIDFromStore(w, req)
	if !ok {
		return
	}
	tpl.RenderTemplate(w, "byID.gohtml", post)
}

func (c BlogController) EditByID(w http.ResponseWriter, req *http.Request) {
	post, ok := c.getPostByIDFromStore(w, req)
	if !ok {
		return
	}

	tpl.RenderTemplate(w, "edit.gohtml", post)
}

// UpdateByID accepts a partial Post and update its fields
func (c BlogController) UpdateByID(w http.ResponseWriter, req *http.Request) {
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
func (c BlogController) DeleteByID(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postID := vars["id"]
	delete(c.store, postID)
	w.Write([]byte("OK"))
}

func (c BlogController) Favicon(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "public/favicon.ico")
}

// LoggingMiddleware logs all incoming requests
// TODO: move to separate controller
func (c BlogController) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("controller.LoggingMiddleware:", req.Method, req.RequestURI)
		next.ServeHTTP(w, req)
	})
}
