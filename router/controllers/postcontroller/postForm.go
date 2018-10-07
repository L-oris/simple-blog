package postcontroller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
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

func parsePostForm(w http.ResponseWriter, req *http.Request, checkRequiredFields bool) (postForm, error) {
	req.Body = http.MaxBytesReader(w, req.Body, maxUploadSize)
	if err := req.ParseMultipartForm(maxUploadSize); err != nil {
		err = fmt.Errorf("uploaded file is too big: %s", err.Error())
		logger.Log.Debug(err.Error())
		return postForm{}, err
	}

	partialPost, err := getPostFromForm(req, checkRequiredFields)
	if err != nil {
		logger.Log.Debug(err.Error())
		return postForm{}, err
	}

	contentType, fileBytes, err := getImageFromForm(req, "postImage")
	if err != nil {
		return postForm{}, err
	}
	if checkRequiredFields && len(fileBytes) == 0 {
		return postForm{}, fmt.Errorf("picture required for post")
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

func getPostFromForm(req *http.Request, checkRequiredFields bool) (post.Post, error) {
	partialPost := post.Post{
		Title:   req.Form["title"][0],
		Content: req.Form["content"][0],
	}

	if checkRequiredFields && !partialPost.HasTitleAndContent() {
		return post.Post{}, errors.New("empty title or content provided")
	}

	return partialPost, nil
}

func getImageFromForm(req *http.Request, inputField string) (contentType string, fileBytes []byte, err error) {
	var multipartFile multipart.File
	multipartFile, _, err = req.FormFile(inputField)
	if err != nil {
		logger.Log.Debug("could not get file from template <form>: %s", err.Error())
		return "", []byte{}, nil
	}
	defer multipartFile.Close()

	fileBytes, err = ioutil.ReadAll(multipartFile)
	if err != nil {
		err = fmt.Errorf("invalid file uploaded: %s", err.Error())
		logger.Log.Debug(err.Error())
		return "", nil, err
	}

	contentType = http.DetectContentType(fileBytes)
	if ok := checkImageType(contentType); !ok {
		err = fmt.Errorf("invalid fileType provided: %s", contentType)
		logger.Log.Debug(err.Error())
		return "", nil, err
	}

	return
}

func checkImageType(fileType string) bool {
	if fileType != "image/jpeg" &&
		fileType != "image/jpg" &&
		fileType != "image/gif" &&
		fileType != "image/png" {
		return false
	}
	return true
}
