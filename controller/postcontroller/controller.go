package postcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/httperror"
	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/post"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
)

type Config struct {
	PathPrefix string
	Tpl        tpl.Template
}

type postControllerStore map[string]post.Post

type postController struct {
	store  postControllerStore
	tpl    tpl.Template
	Router *mux.Router
}

// New creates a new controller and registers the routes
func New(config *Config) postController {
	c := postController{
		store: make(map[string]post.Post),
		tpl:   config.Tpl,
	}

	router := mux.NewRouter()
	routes := router.PathPrefix(config.PathPrefix).Subrouter()
	routes.HandleFunc("/all", c.getAll).Methods("GET")
	routes.HandleFunc("/new", c.new).Methods("GET")
	routes.HandleFunc("/new", c.add).Methods("POST")
	routes.HandleFunc("/{id}", c.getByID).Methods("GET")
	routes.HandleFunc("/{id}/edit", c.editByID).Methods("GET")
	routes.HandleFunc("/{id}/edit", c.updateByID).Methods("POST")
	routes.HandleFunc("/{id}", c.deleteByID).Methods("DELETE")

	c.Router = router
	return c
}

// getAll gets all Post from the store
func (c postController) getAll(w http.ResponseWriter, req *http.Request) {
	c.tpl.Render(w, "all.gohtml", c.store)
}

// new renders the template for adding new Post
func (c postController) new(w http.ResponseWriter, req *http.Request) {
	c.tpl.Render(w, "new.gohtml", nil)
}

// add adds a Post to the store
func (c postController) add(w http.ResponseWriter, req *http.Request) {
	partialPost, err := getPartialPostFromForm(req, true)
	if err != nil {
		logger.Log.Warning("incomplete post received")
		httperror.BadRequest(w, err.Error())
		return
	}

	newPost, _ := post.GenerateFromPartial(partialPost)
	c.store[newPost.ID] = newPost
	c.tpl.Render(w, "byID.gohtml", newPost)
}

// getByID gets a Post by ID from store
func (c postController) getByID(w http.ResponseWriter, req *http.Request) {
	post, ok := c.getPostByIDFromStore(w, req)
	if !ok {
		httperror.NotFound(w, "Post not found")
		return
	}
	c.tpl.Render(w, "byID.gohtml", post)
}

func (c postController) editByID(w http.ResponseWriter, req *http.Request) {
	post, ok := c.getPostByIDFromStore(w, req)
	if !ok {
		httperror.NotFound(w, "Post not found")
		return
	}

	c.tpl.Render(w, "edit.gohtml", post)
}

// updateByID accepts a partial Post and update its fields
func (c postController) updateByID(w http.ResponseWriter, req *http.Request) {
	storePost, ok := c.getPostByIDFromStore(w, req)
	if !ok {
		httperror.NotFound(w, "Post not found")
		return
	}

	newPartialPost, err := getPartialPostFromForm(req, false)
	if err != nil {
		logger.Log.Warning("incomplete post received")
		httperror.BadRequest(w, err.Error())
		return
	}

	if err = mergo.Merge(&storePost, newPartialPost, mergo.WithOverride); err != nil {
		logger.Log.Error("failed to merge posts")
		httperror.BadRequest(w, "Invalid post provided")
		return
	}

	c.store[storePost.ID] = storePost
	c.tpl.Render(w, "byID.gohtml", storePost)
}

// deleteByID deletes a Post by ID
// It doesn't care if the Post exists or not
func (c postController) deleteByID(w http.ResponseWriter, req *http.Request) {
	postID := mux.Vars(req)["id"]
	delete(c.store, postID)
	w.Write([]byte("OK"))
}
