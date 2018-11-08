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
	"github.com/L-oris/yabb/router/controllers/postcontroller"
	"github.com/L-oris/yabb/router/controllers/rootcontroller"
	"github.com/L-oris/yabb/services/postservice"
)

func FileServer() (func(w http.ResponseWriter, r *http.Request, name string), error) {
	return http.ServeFile, nil
}

func BlogDB() (*sql.DB, error) {
	return db.BlogDB, nil
}

func Templates() (*template.Template, error) {
	return &template.Template{}, nil
}

func PostRepository(db *sql.DB) (*postrepository.Repository, error) {
	return postrepository.New(
		&postrepository.Config{
			DB: db,
		},
	), nil
}

func BucketRepository() (*bucketrepository.Repository, error) {
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

func PostService(bucket *bucketrepository.Repository, repository *postrepository.Repository) (*postservice.Service, error) {
	return postservice.New(&postservice.Config{
		Bucket:     bucket,
		Repository: repository,
	}), nil
}

func RootController(renderer *template.Template, serve func(w http.ResponseWriter, r *http.Request, fileName string), bucket *bucketrepository.Repository, db *sql.DB) (rootcontroller.Controller, error) {
	return rootcontroller.New(
		&rootcontroller.Config{
			Renderer: renderer,
			Serve:    serve,
			Bucket:   bucket,
			DB:       db,
		}), nil
}

func PostController(renderer *template.Template, service *postservice.Service) (postcontroller.Controller, error) {
	return postcontroller.New(&postcontroller.Config{
		Renderer: renderer,
		Service:  service,
	}), nil
}
