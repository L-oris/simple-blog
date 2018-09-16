package postcontroller

import (
	"errors"
	"net/http"

	"github.com/L-oris/yabb/httperror"
	"github.com/L-oris/yabb/models/post"
	"github.com/gorilla/mux"
)

// getPostByIDFromStore gets ID from url params and returns Post
// When not found, returns false as second param
func (c postController) getPostByIDFromStore(w http.ResponseWriter, req *http.Request) (post.Post, bool) {
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

// getPartialPostFromForm parses request form and returns a post with Title & Content (other values are zeroed)
// 'checkTitleAndContent' defines whether title & content should be mandatory
func getPartialPostFromForm(req *http.Request, checkTitleAndContent bool) (post.Post, error) {
	req.ParseForm()

	if len(req.Form["title"]) == 0 || len(req.Form["content"]) == 0 {
		return post.Post{}, errors.New("Invalid data provided")
	}

	partialPost := post.Post{
		Title:   req.Form["title"][0],
		Content: req.Form["content"][0],
	}

	if checkTitleAndContent && !partialPost.HasTitleAndContent() {
		return post.Post{}, errors.New("Empty title or content provided")
	}

	return partialPost, nil
}
