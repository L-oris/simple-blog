package postservice

import (
	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/repository/bucketrepository"
	"github.com/L-oris/yabb/repository/postrepository"
)

type Config struct {
	Bucket     *bucketrepository.Repository
	Repository *postrepository.Repository
}

type Service struct {
	bucket     *bucketrepository.Repository
	repository *postrepository.Repository
}

// New creates a new PostService
func New(config *Config) *Service {
	return &Service{
		bucket:     config.Bucket,
		repository: config.Repository,
	}
}

// Create creates a new post
func (s Service) Create() {
	logger.Log.Info("workds!")
}
