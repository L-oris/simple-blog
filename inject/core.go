package inject

import (
	"net/http"

	"github.com/L-oris/yabb/models/tpl"
	"github.com/sarulabs/di"
)

func getCore() []di.Def {
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

	return []di.Def{
		templates, fileserver,
	}
}
