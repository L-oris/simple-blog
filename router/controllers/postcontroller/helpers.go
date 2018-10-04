package postcontroller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/post"
	"github.com/gorilla/mux"
)

// getPostIDFromURL gets 'id' from url query params
func getPostIDFromURL(req *http.Request) (int, error) {
	vars := mux.Vars(req)
	pID, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.Log.Warning("bad id received" + string(pID))
		return 0, err
	}

	return pID, nil
}

// getPartialPostFromForm parses request form and returns a post with Title & Content (other values are zeroed)
// 'checkTitleAndContent' param defines whether title & content should be mandatory
func getPartialPostFromForm(req *http.Request, checkTitleAndContent bool) (post.Post, error) {
	partialPost := post.Post{
		Title:   req.Form["title"][0],
		Content: req.Form["content"][0],
	}

	if checkTitleAndContent && !partialPost.HasTitleAndContent() {
		logger.Log.Warning("post missing Title or Content")
		return post.Post{}, errors.New("empty title or content provided")
	}

	return partialPost, nil
}

func parseMultiPartImage(req *http.Request, inputField string) (contentType string, fileBytes []byte, err error) {
	var multipartFile multipart.File
	multipartFile, _, err = req.FormFile(inputField)
	if err != nil {
		logger.Log.Error("could not get form from template: %s", err.Error())
		return "", nil, err
	}
	defer multipartFile.Close()

	fileBytes, err = ioutil.ReadAll(multipartFile)
	if err != nil {
		logger.Log.Debug("invalid file uploaded: %s", err.Error())
		return "", nil, err
	}

	contentType = http.DetectContentType(fileBytes)
	if ok := checkImageType(contentType); !ok {
		return "", nil, fmt.Errorf("invalid fileType provided: %s", err.Error())
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
