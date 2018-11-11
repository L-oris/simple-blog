//+build wireinject

package injector

import (
	"github.com/L-oris/yabb/services/postservice"
	"github.com/google/go-cloud/wire"
)

func providePostService() (*postservice.Service, error) {
	wire.Build(postservice.New, provideBucket, providePostRepository)
	return &postservice.Service{}, nil
}
