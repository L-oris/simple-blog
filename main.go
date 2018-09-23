package main

import (
	"net/http"
	"time"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/env"
	"github.com/L-oris/yabb/router"
)

func main() {
	server := &http.Server{
		Addr:         ":" + env.Vars.Port,
		Handler:      router.Mount(),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	logger.Log.Fatal(server.ListenAndServe())
}
