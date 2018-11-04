package inject

import (
	"github.com/L-oris/yabb/inject/types"
	"github.com/L-oris/yabb/repository/bucketrepository"
	"github.com/L-oris/yabb/repository/postrepository"
	"github.com/L-oris/yabb/services/postservice"
	"github.com/sarulabs/di"
)

func services() []di.Def {
	postServiceValue := di.Def{
		Name: types.PostService.String(),
		Build: func(ctn di.Container) (interface{}, error) {
			return postservice.New(&postservice.Config{
				Bucket:     ctn.Get(types.BucketRepository.String()).(*bucketrepository.Repository),
				Repository: ctn.Get(types.PostRepository.String()).(*postrepository.Repository),
			}), nil
		},
	}

	return []di.Def{postServiceValue}
}
