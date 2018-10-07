package postcontroller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/L-oris/yabb/logger"
	"github.com/gorilla/mux"
)

// getPostIDFromURL gets 'id' from url query params
func getPostIDFromURL(req *http.Request) (int, error) {
	vars := mux.Vars(req)
	pID, err := strconv.Atoi(vars["id"])
	if err != nil {
		err = fmt.Errorf("bad id received" + string(pID))
		logger.Log.Warning(err.Error())
		return 0, err
	}

	return pID, nil
}
