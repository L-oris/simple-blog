package inject

import (
	"net/http"

	"github.com/L-oris/yabb/models/tpl"
	"github.com/sarulabs/di"
)

// Container stores all dependencies, allowing to easily inject them
var Container di.Container

func init() {
	Container = createBuilder().Build()
}

func createBuilder() *di.Builder {
	tpl := di.Def{
		Name: "tpl",
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
	builder.Add(append(createRepositories(), fileserver, tpl)...)

	return builder
}
