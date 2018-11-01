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

const maxUploadedImageSize = 1 * 1024 * 1024 // MB

// postForm contains data being parsed in template form submission
type postForm struct {
	post      post.Post
	fileBytes []byte
}

func parsePostForm(w http.ResponseWriter, req *http.Request, checkRequiredFields bool) (postForm, error) {
	req.Body = http.MaxBytesReader(w, req.Body, maxUploadedImageSize)
	if err := req.ParseMultipartForm(maxUploadedImageSize); err != nil {
		return postForm{}, errors.New("uploaded file is too big")
	}

	post := getTextFieldsFromForm(req)
	if checkRequiredFields && !post.HasTitleAndContent() {
		return postForm{}, errors.New("missing mandatory fields for post")
	}

	contentType, fileBytes, err := getImageFromForm(req, "postImage")
	if err != nil {
		return postForm{}, err
	}
	if len(fileBytes) == 0 {
		if checkRequiredFields {
			return postForm{}, errors.New("missing mandatory picture for post")
		}

		return postForm{
			post:      post,
			fileBytes: nil,
		}, nil
	}

	fileEndings, _ := mime.ExtensionsByType(contentType)
	fileName := uuid.NewV4().String() + fileEndings[0]
	post.Picture = fileName

	logger.Log.Debug("ContentType: %s, File: %s", contentType, fileName)
	return postForm{
		post:      post,
		fileBytes: fileBytes,
	}, nil
}

func getTextFieldsFromForm(req *http.Request) post.Post {
	return post.Post{
		Title:   req.Form["title"][0],
		Content: req.Form["content"][0],
	}
}

func getImageFromForm(req *http.Request, inputField string) (contentType string, fileBytes []byte, err error) {
	var multipartFile multipart.File
	multipartFile, _, err = req.FormFile(inputField)
	if err != nil {
		logger.Log.Debug("no image submitted, proceeding with empty file")
		return "", []byte{}, nil
	}
	defer multipartFile.Close()

	fileBytes, err = ioutil.ReadAll(multipartFile)
	if err != nil {
		logger.Log.Debugf("invalid file uploaded: %s", err.Error())
		return "", nil, errors.New("invalid file uploaded")
	}

	contentType = http.DetectContentType(fileBytes)
	if ok := checkImageType(contentType); !ok {
		return "", nil, fmt.Errorf("invalid fileType provided: %s", contentType)
	}

	return contentType, fileBytes, nil
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
