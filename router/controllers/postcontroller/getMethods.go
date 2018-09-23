package postcontroller

import (
	"net/http"

	"github.com/L-oris/yabb/router/httperror"
)

func (c Controller) renderAll(w http.ResponseWriter, req *http.Request) {
	posts, err := c.repository.GetAll()
	if err != nil {
		httperror.BadRequest(w, "cannot get posts")
	}
	c.tpl.Render(w, "all.gohtml", posts)
}

func (c Controller) renderNew(w http.ResponseWriter, req *http.Request) {
	c.tpl.Render(w, "new.gohtml", nil)
}

func (c Controller) renderByID(w http.ResponseWriter, req *http.Request) {
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

func (c Controller) renderUpdateByID(w http.ResponseWriter, req *http.Request) {
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
