package inject

import (
	"github.com/L-oris/yabb/inject/types"
	"github.com/L-oris/yabb/router"
	"github.com/sarulabs/di"
)

func routers() []di.Def {
	routerValue := di.Def{
		Name: types.Router.String(),
		Build: func(ctn di.Container) (interface{}, error) {
			return router.Mount(ctn), nil
		},
	}

	return []di.Def{
		routerValue,
	}
}
