package main

import (
	"net/http"
	"time"

	"github.com/L-oris/yabb/inject"
	"github.com/L-oris/yabb/inject/types"
	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/env"
)

func main() {
	container := inject.CreateContainer()
	router := container.Get(types.Router.String()).(http.Handler)
	server := &http.Server{
		Addr:         ":" + env.Vars.Port,
		Handler:      router,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	logger.Log.Infof("server listening on port %s", env.Vars.Port)
	logger.Log.Fatal(server.ListenAndServe())
}

// TODO: move 'tpl' and 'env' packages out of 'models' directory
// TODO: move away post actions from postRepository (create separate engine)
// TODO: 'httperror' package should render error template with error message we receive
