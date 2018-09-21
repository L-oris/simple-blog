package inject

import (
	"net/http"

	"github.com/L-oris/yabb/controller/postcontroller"
	"github.com/L-oris/yabb/controller/rootcontroller"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/L-oris/yabb/repository/postrepository"
	"github.com/sarulabs/di"
)

func getControllers() []di.Def {
	rootControllerValue := di.Def{
		Name: "rootcontroller",
		Build: func(ctn di.Container) (interface{}, error) {
			return rootcontroller.New(
				&rootcontroller.Config{
					PathPrefix: "/",
					Tpl:        ctn.Get("templates").(*tpl.TPL),
					Serve:      ctn.Get("fileserver").(func(w http.ResponseWriter, r *http.Request, fileName string)),
				}), nil
		},
	}

	postControllerValue := di.Def{
		Name: "postcontroller",
		Build: func(ctn di.Container) (interface{}, error) {
			return postcontroller.New(&postcontroller.Config{
				PathPrefix: "/post",
				Repository: ctn.Get("postrepository").(*postrepository.Repository),
				Tpl:        ctn.Get("templates").(*tpl.TPL),
			}), nil
		},
	}

	return []di.Def{
		rootControllerValue, postControllerValue,
	}
}
