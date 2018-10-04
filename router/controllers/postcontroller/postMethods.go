package postcontroller

import (
	"mime"
	"net/http"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/router/httperror"
	uuid "github.com/satori/go.uuid"
)

const maxUploadSize = 2 * 1024 * 1024 // MB

func (c Controller) new(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, maxUploadSize)
	if err := req.ParseMultipartForm(maxUploadSize); err != nil {
		logger.Log.Debug("uploaded file is too big: %s", err.Error())
		httperror.BadRequest(w, "file provided is too large")
		return
	}

	contentType, fileBytes, err := parseMultiPartImage(req, "postImage")
	if err != nil {
		httperror.BadRequest(w, "invalid file provided")
	}

	fileEndings, _ := mime.ExtensionsByType(contentType)
	fileName := uuid.NewV4().String() + fileEndings[0]
	logger.Log.Debug("ContentType: %s, File: %s", contentType, fileName)

	partialPost, err := getPartialPostFromForm(req, true)
	if err != nil {
		httperror.BadRequest(w, "incomplete post received")
		return
	}

	partialPost.Picture = fileName
	newPost, err := c.repository.Add(partialPost)
	if err != nil {
		httperror.InternalServer(w, "failed to save post")
		return
	}

	err = c.bucket.Write(fileName, fileBytes)
	if err != nil {
		httperror.InternalServer(w, "cannot save file")
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

// TODO:
// * delete from bucket when deleting post
