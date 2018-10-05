package postcontroller

import (
	"fmt"
	"mime"
	"net/http"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/post"
	uuid "github.com/satori/go.uuid"
)

const maxUploadSize = 1 * 1024 * 1024 // MB

// postForm contains data being parsed in template form submission
type postForm struct {
	post      post.Post
	fileBytes []byte
}

func parsePostForm(w http.ResponseWriter, req *http.Request) (postForm, error) {
	req.Body = http.MaxBytesReader(w, req.Body, maxUploadSize)
	if err := req.ParseMultipartForm(maxUploadSize); err != nil {
		err = fmt.Errorf("uploaded file is too big: %s", err.Error())
		logger.Log.Debug(err.Error())
		return postForm{}, err
	}

	partialPost, err := getPostFromForm(req, true)
	if err != nil {
		err = fmt.Errorf("incomplete post received")
		logger.Log.Debug(err.Error())
		return postForm{}, err
	}

	contentType, fileBytes, err := getImageFromForm(req, "postImage")
	if err != nil {
		if err != nil {
			err = fmt.Errorf("invalid file type provided: %s", err.Error())
			logger.Log.Debug(err.Error())
			return postForm{}, err
		}
	}
	fileEndings, _ := mime.ExtensionsByType(contentType)
	fileName := uuid.NewV4().String() + fileEndings[0]
	partialPost.Picture = fileName
	logger.Log.Debug("ContentType: %s, File: %s", contentType, fileName)

	return postForm{
		post:      partialPost,
		fileBytes: fileBytes,
	}, nil
}
