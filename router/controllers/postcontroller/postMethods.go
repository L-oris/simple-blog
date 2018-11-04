package postcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/router/httperror"
)

func (c Controller) new(w http.ResponseWriter, req *http.Request) {
	postForm, err := parsePostForm(w, req, true)
	if err != nil {
		httperror.BadRequest(w, err.Error())
		return
	}

	newPost, err := c.service.Create(postForm.post, postForm.fileBytes)
	if err != nil {
		httperror.InternalServer(w, err.Error())
		return
	}

	w.Header().Set("Location", "/post/"+newPost.ID)
	w.WriteHeader(http.StatusSeeOther)
}

func (c Controller) updateByID(w http.ResponseWriter, req *http.Request) {
	pID, err := getPostIDFromURL(req)
	if err != nil {
		httperror.BadRequest(w, err.Error())
		return
	}

	oldPost, err := c.repository.GetByID(pID)
	if err != nil {
		httperror.BadRequest(w, err.Error())
		return
	}

	postForm, err := parsePostForm(w, req, false)
	if err != nil {
		httperror.BadRequest(w, err.Error())
		return
	}

	post, err := c.repository.UpdateByID(pID, postForm.post)
	if err != nil {
		httperror.InternalServer(w, err.Error())
		return
	}

	if len(postForm.fileBytes) > 0 {
		if err = c.bucket.Write(postForm.post.Picture, postForm.fileBytes); err != nil {
			httperror.InternalServer(w, err.Error())
			return
		}

		if err = c.bucket.Delete(oldPost.Picture); err != nil {
			httperror.InternalServer(w, err.Error())
			return
		}
	}

	w.Header().Set("Location", "/post/"+post.ID)
	w.WriteHeader(http.StatusSeeOther)
}

func (c Controller) deleteByID(w http.ResponseWriter, req *http.Request) {
	pID, err := getPostIDFromURL(req)
	if err != nil {
		httperror.BadRequest(w, err.Error())
		return
	}

	post, err := c.repository.GetByID(pID)
	if err != nil {
		httperror.BadRequest(w, err.Error())
		return
	}

	if err = c.bucket.Delete(post.Picture); err != nil {
		httperror.InternalServer(w, err.Error())
		return
	}

	if err = c.repository.DeleteByID(pID); err != nil {
		httperror.InternalServer(w, err.Error())
		return
	}

	w.Header().Set("Location", "/post/all")
	w.WriteHeader(http.StatusSeeOther)
}
