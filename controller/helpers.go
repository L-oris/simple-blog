package controller

import (
	"net/http"

	"github.com/L-oris/yabb/httperror"
	"github.com/L-oris/yabb/models/post"
	"github.com/gorilla/mux"
)

// getPostByIDFromStore gets ID from url params and returns Post
// When not found, returns false as second param
func (c BlogController) getPostByIDFromStore(w http.ResponseWriter, req *http.Request) (post.Post, bool) {
	vars := mux.Vars(req)
	pID := vars["id"]
	p, ok := c.store[pID]
	if !ok {
		errMessage := "Post " + pID + " not found"
		httperror.NotFound(w, errMessage)
		return post.Post{}, false
	}

	return p, true
}

// func getFromStore(w http.ResponseWriter, req *http.Request) post.Post {}
