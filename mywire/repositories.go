//+build wireinject

package mywire

import (
	"fmt"

	"github.com/L-oris/yabb/foreign/env"
	"github.com/L-oris/yabb/repositories/bucketrepository"
	"github.com/L-oris/yabb/repositories/postrepository"
	"github.com/google/go-cloud/wire"
)

func provideBucket() (*bucketrepository.Repository, error) {
	repo, err := bucketrepository.New(bucketrepository.BucketName(env.Vars.BucketName))
	if err != nil {
		return nil, fmt.Errorf("could not create bucket: %s", err.Error())
	}
	return repo, nil
}
func providePostRepository() (*postrepository.Repository, error) {
	wire.Build(postrepository.New, provideDB)
	return &postrepository.Repository{}, nil
}
