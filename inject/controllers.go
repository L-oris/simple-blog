package inject

import (
	"net/http"

	"github.com/L-oris/yabb/repository/bucketrepository"

	"github.com/L-oris/yabb/inject/types"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/L-oris/yabb/repository/postrepository"
	"github.com/L-oris/yabb/router/controllers/postcontroller"
	"github.com/L-oris/yabb/router/controllers/rootcontroller"
	"github.com/sarulabs/di"
)

func controllers() []di.Def {
	rootControllerValue := di.Def{
		Name: types.RootController.String(),
		Build: func(ctn di.Container) (interface{}, error) {
			return rootcontroller.New(
				&rootcontroller.Config{
					Tpl:    ctn.Get(types.Templates.String()).(*tpl.TPL),
					Serve:  ctn.Get(types.FileServer.String()).(func(w http.ResponseWriter, r *http.Request, fileName string)),
					Bucket: ctn.Get(types.BucketRepository.String()).(*bucketrepository.Repository),
				}), nil
		},
	}

	postControllerValue := di.Def{
		Name: types.PostController.String(),
		Build: func(ctn di.Container) (interface{}, error) {
			return postcontroller.New(&postcontroller.Config{
				Repository: ctn.Get(types.PostRepository.String()).(*postrepository.Repository),
				Tpl:        ctn.Get(types.Templates.String()).(*tpl.TPL),
			}), nil
		},
	}

	return []di.Def{
		rootControllerValue, postControllerValue,
	}
}
