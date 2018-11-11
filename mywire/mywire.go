//+build wireinject

package mywire

import (
	"net/http"

	"github.com/L-oris/yabb/router"
	"github.com/google/go-cloud/wire"
)

func InitializeRouter() (http.Handler, error) {
	wire.Build(router.New, router.Config{}, provideRootController, providePostController)
	return nil, nil
}
