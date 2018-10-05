package postcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/router/httperror"
)

func (c Controller) new(w http.ResponseWriter, req *http.Request) {
	postForm, err := parsePostForm(w, req)
	if err != nil {
		httperror.BadRequest(w, err.Error())
		return
	}

	newPost, err := c.repository.Add(postForm.post)
	if err != nil {
		httperror.InternalServer(w, "failed to save post")
		return
	}

	if err = c.bucket.Write(postForm.post.Picture, postForm.fileBytes); err != nil {
		httperror.InternalServer(w, "failed to save file")
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

	partialPost, err := getPostFromForm(req, false)
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

// TODO:
// * delete from bucket when deleting post
