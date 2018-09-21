package inject

import (
	"net/http"

	"github.com/L-oris/yabb/controller/rootcontroller"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/sarulabs/di"
)

// Container stores all dependencies, allowing to easily inject them
var Container di.Container

func init() {
	Container = createBuilder().Build()
}

func createBuilder() *di.Builder {
	templates := di.Def{
		Name: "templates",
		Build: func(ctn di.Container) (interface{}, error) {
			return &tpl.TPL{}, nil
		},
	}

	fileserver := di.Def{
		Name: "fileserver",
		Build: func(ctn di.Container) (interface{}, error) {
			return http.ServeFile, nil
		},
	}

	builder, _ := di.NewBuilder()
	builder.Add(fileserver, templates)
	builder.Add(createRepositories()...)
	builder.Add(createControllers()...)

	return builder
}

func createControllers() []di.Def {
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

	return []di.Def{
		rootControllerValue,
	}
}
