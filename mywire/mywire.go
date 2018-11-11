//+build wireinject

package mywire

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/L-oris/yabb/foreign/env"
	"github.com/L-oris/yabb/foreign/template"
	"github.com/L-oris/yabb/repositories/bucketrepository"
	"github.com/L-oris/yabb/repositories/db"
	"github.com/L-oris/yabb/repositories/postrepository"
	"github.com/L-oris/yabb/router"
	"github.com/L-oris/yabb/router/controllers/postcontroller"
	"github.com/L-oris/yabb/router/controllers/rootcontroller"
	"github.com/L-oris/yabb/services/postservice"
	"github.com/google/go-cloud/wire"
)

func provideFileServer() (func(w http.ResponseWriter, r *http.Request, name string), error) {
	return http.ServeFile, nil
}

func provideDB() *sql.DB {
	return db.BlogDB
}

func provideRenderer() (template.Renderer, error) {
	return template.Template{}, nil
}

func provideBucket() (*bucketrepository.Repository, error) {
	repo, err := bucketrepository.New(bucketrepository.BucketName(env.Vars.BucketName))
	if err != nil {
		return nil, fmt.Errorf("could not create bucket: %s", err.Error())
	}
	return repo, nil
}

func provideRootController() (rootcontroller.Controller, error) {
	wire.Build(rootcontroller.New, rootcontroller.Config{}, provideFileServer, provideDB, provideRenderer, provideBucket)
	return rootcontroller.Controller{}, nil
}

func providePostRepository() (*postrepository.Repository, error) {
	wire.Build(postrepository.New, provideDB)
	return &postrepository.Repository{}, nil
}

func providePostService() (*postservice.Service, error) {
	wire.Build(postservice.New, provideBucket, providePostRepository)
	return &postservice.Service{}, nil
}

func providePostController() (postcontroller.Controller, error) {
	wire.Build(postcontroller.New, postcontroller.Config{}, provideRenderer, providePostService)
	return postcontroller.Controller{}, nil
}

func InitializeRouter() (http.Handler, error) {
	wire.Build(router.New, router.Config{}, provideRootController, providePostController)
	return nil, nil
}
