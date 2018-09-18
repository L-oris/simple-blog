package postcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/httperror"
	"github.com/L-oris/yabb/models/post"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/L-oris/yabb/repository/postrepository"
	"github.com/gorilla/mux"
)

type Config struct {
	PathPrefix string
	Repository *postrepository.Repository
	Tpl        tpl.Template
}

type postControllerStore map[string]post.Post

type postController struct {
	repository *postrepository.Repository
	store      postControllerStore
	tpl        tpl.Template
	Router     *mux.Router
}

// New creates a new controller and registers the routes
func New(config *Config) postController {
	c := postController{
		repository: config.Repository,
		store:      make(map[string]post.Post),
		tpl:        config.Tpl,
	}

	router := mux.NewRouter()
	routes := router.PathPrefix(config.PathPrefix).Subrouter()
	routes.HandleFunc("/ping", c.ping).Methods("GET")
	routes.HandleFunc("/all", c.renderAll).Methods("GET")
	routes.HandleFunc("/new", c.renderNew).Methods("GET")
	routes.HandleFunc("/{id}", c.renderByID).Methods("GET")
	routes.HandleFunc("/{id}/update", c.renderUpdateByID).Methods("GET")
	routes.HandleFunc("/new", c.new).Methods("POST")
	routes.HandleFunc("/{id}/update", c.updateByID).Methods("POST")
	routes.HandleFunc("/{id}/delete", c.deleteByID).Methods("POST")

	c.Router = router
	return c
}

// ping checks db connection
func (c postController) ping(w http.ResponseWriter, req *http.Request) {
	if err := c.repository.Ping(); err != nil {
		httperror.InternalServer(w, "db not connected")
		return
	}
	w.Write([]byte("pong"))
}

func (c postController) renderAll(w http.ResponseWriter, req *http.Request) {
	posts, err := c.repository.GetAll()
	if err != nil {
		httperror.BadRequest(w, "cannot get posts")
	}
	c.tpl.Render(w, "all.gohtml", posts)
}

func (c postController) renderNew(w http.ResponseWriter, req *http.Request) {
	c.tpl.Render(w, "new.gohtml", nil)
}

func (c postController) new(w http.ResponseWriter, req *http.Request) {
	partialPost, err := getPartialPostFromForm(req, true)
	if err != nil {
		httperror.BadRequest(w, "incomplete post received")
		return
	}

	newPost, err := c.repository.Add(partialPost)
	if err != nil {
		httperror.InternalServer(w, "failed to add post")
		return
	}

	w.Header().Set("Location", "/post/"+newPost.ID)
	w.WriteHeader(http.StatusSeeOther)
}

func (c postController) renderByID(w http.ResponseWriter, req *http.Request) {
	pID, err := getPostIDFromURL(req)
	if err != nil {
		httperror.BadRequest(w, "bad id provided: "+string(pID))
	}

	post, err := c.repository.GetByID(pID)
	if err != nil {
		httperror.NotFound(w, "Post "+string(pID)+" not found")
	}

	c.tpl.Render(w, "byID.gohtml", post)
}

func (c postController) renderUpdateByID(w http.ResponseWriter, req *http.Request) {
	pID, err := getPostIDFromURL(req)
	if err != nil {
		httperror.BadRequest(w, "bad id provided: "+string(pID))
	}

	post, err := c.repository.GetByID(pID)
	if err != nil {
		httperror.NotFound(w, "Post "+string(pID)+" not found")
	}

	c.tpl.Render(w, "edit.gohtml", post)
}

func (c postController) updateByID(w http.ResponseWriter, req *http.Request) {
	pID, err := getPostIDFromURL(req)
	if err != nil {
		httperror.BadRequest(w, "bad id provided: "+string(pID))
	}

	partialPost, err := getPartialPostFromForm(req, false)
	if err != nil {
		httperror.BadRequest(w, "invalid data provided")
		return
	}

	post, err := c.repository.UpdateByID(pID, partialPost)

	w.Header().Set("Location", "/post/"+post.ID)
	w.WriteHeader(http.StatusSeeOther)
}

func (c postController) deleteByID(w http.ResponseWriter, req *http.Request) {
	pID, err := getPostIDFromURL(req)
	if err != nil {
		httperror.BadRequest(w, "bad id provided: "+string(pID))
	}

	if err = c.repository.DeleteByID(pID); err != nil {
		httperror.BadRequest(w, "cannot delete post "+string(pID))
		return
	}

	w.Header().Set("Location", "/post/all")
	w.WriteHeader(http.StatusSeeOther)
}
