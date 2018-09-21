package inject

import (
	"github.com/L-oris/yabb/models/db"
	"github.com/L-oris/yabb/repository/postrepository"
	"github.com/sarulabs/di"
)

func getRepositories() []di.Def {
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

	return []di.Def{
		postRepository,
	}
}
