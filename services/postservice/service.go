package postservice

import (
	"github.com/L-oris/yabb/models/post"
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
func (s Service) Create(newPost post.Post, fileBytes []byte) (post.Post, error) {
	dbPost, err := s.repository.Add(newPost)
	if err != nil {
		return post.Post{}, err
	}

	if err = s.bucket.Write(dbPost.Picture, fileBytes); err != nil {
		return post.Post{}, err
	}

	return dbPost, nil
}
