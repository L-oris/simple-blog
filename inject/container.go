package inject

import (
	"net/http"

	"github.com/L-oris/yabb/models/db"
	"github.com/L-oris/yabb/models/tpl"
	"github.com/L-oris/yabb/repository/postrepository"
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

	postRepository := di.Def{
		Name: "postrepository",
		Build: func(ctn di.Container) (interface{}, error) {
			return postrepository.New(
				&postrepository.Config{
					DB: db.BlogDB,
				},
			), nil
		},
	}

	fileserver := di.Def{
		Name: "fileserver",
		Build: func(ctn di.Container) (interface{}, error) {
			return http.ServeFile, nil
		},
	}

	builder, _ := di.NewBuilder()
	builder.Add(tpl, postRepository, fileserver)

	return builder
}
