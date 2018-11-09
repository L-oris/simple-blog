package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/L-oris/yabb/foreign/env"
	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/mywire"
)

func main() {
	router, err := mywire.ProvideRouter()
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(1)
	}
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
