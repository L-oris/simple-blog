package postcontroller

import (
	"errors"
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
	req.ParseForm()
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
