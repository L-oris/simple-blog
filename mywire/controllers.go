//+build wireinject

package mywire

import (
	"github.com/L-oris/yabb/router/controllers/postcontroller"
	"github.com/L-oris/yabb/router/controllers/rootcontroller"
	"github.com/google/go-cloud/wire"
)

func provideRootController() (rootcontroller.Controller, error) {
	wire.Build(rootcontroller.New, rootcontroller.Config{}, provideFileServer, provideDB, provideRenderer, provideBucket)
	return rootcontroller.Controller{}, nil
}

func providePostController() (postcontroller.Controller, error) {
	wire.Build(postcontroller.New, postcontroller.Config{}, provideRenderer, providePostService)
	return postcontroller.Controller{}, nil
}
