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

func ProvideFileServer() (func(w http.ResponseWriter, r *http.Request, name string), error) {
	return http.ServeFile, nil
}

func ProvideBlogDB() *sql.DB {
	return db.BlogDB
}

func ProvideTemplates() (*template.Template, error) {
	return &template.Template{}, nil
}

func ProvideBucket() (*bucketrepository.Repository, error) {
	repo, err := bucketrepository.New(
		bucketrepository.Config{
			BucketName: env.Vars.BucketName,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not create bucket: %s", err.Error())
	}
	return repo, nil
}

func ProvideRootController() (rootcontroller.Controller, error) {
	wire.Build(rootcontroller.NewWire, ProvideFileServer, ProvideBlogDB, ProvideTemplates, ProvideBucket)
	return rootcontroller.Controller{}, nil
}

func ProvidePostRepository() (*postrepository.Repository, error) {
	wire.Build(postrepository.NewWire, ProvideBlogDB)
	return &postrepository.Repository{}, nil
}

func ProvidePostService() (*postservice.Service, error) {
	wire.Build(postservice.NewWire, ProvideBucket, ProvidePostRepository)
	return &postservice.Service{}, nil
}

func ProvidePostController() (postcontroller.Controller, error) {
	wire.Build(postcontroller.NewWire, ProvideTemplates, ProvidePostService)
	return postcontroller.Controller{}, nil
}

func ProvideRouter() (http.Handler, error) {
	wire.Build(router.NewWire, ProvideRootController, ProvidePostController)
	return nil, nil
}
