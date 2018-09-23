package postcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/router/httperror"
)

func (c Controller) new(w http.ResponseWriter, req *http.Request) {
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

func (c Controller) updateByID(w http.ResponseWriter, req *http.Request) {
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

func (c Controller) deleteByID(w http.ResponseWriter, req *http.Request) {
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
