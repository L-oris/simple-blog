package postcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/router/httperror"
)

func (c Controller) renderAll(w http.ResponseWriter, req *http.Request) {
	posts, err := c.service.GetAll()
	if err != nil {
		httperror.InternalServer(w, err.Error())
		return
	}
	c.renderer.Render(w, "all.gohtml", posts)
}

func (c Controller) renderNew(w http.ResponseWriter, req *http.Request) {
	c.renderer.Render(w, "new.gohtml", nil)
}

func (c Controller) renderByID(w http.ResponseWriter, req *http.Request) {
	postID, err := getPostIDFromURL(req)
	if err != nil {
		httperror.BadRequest(w, err.Error())
		return
	}

	post, err := c.service.GetByID(postID)
	if err != nil {
		httperror.NotFound(w, err.Error())
		return
	}

	c.renderer.Render(w, "byID.gohtml", post)
}

func (c Controller) renderUpdateByID(w http.ResponseWriter, req *http.Request) {
	postID, err := getPostIDFromURL(req)
	if err != nil {
		httperror.BadRequest(w, err.Error())
		return
	}

	post, err := c.service.GetByID(postID)
	if err != nil {
		httperror.NotFound(w, err.Error())
		return
	}

	c.renderer.Render(w, "edit.gohtml", post)
}
