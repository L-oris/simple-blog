package postservice

import (
	"fmt"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/post"
	"github.com/L-oris/yabb/repositories/bucketrepository"
	"github.com/L-oris/yabb/repositories/postrepository"
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

// GetAll gets all posts
func (s Service) GetAll() ([]post.Post, error) {
	dbPosts, err := s.repository.GetAll()
	if err != nil {
		return []post.Post{}, err
	}
	return dbPosts, nil
}

// GetByID gets a post by ID
func (s Service) GetByID(postID int) (post.Post, error) {
	dbPost, err := s.repository.GetByID(postID)
	if err != nil {
		return post.Post{}, err
	}
	return dbPost, nil
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

// UpdateByID updates a post by ID
func (s Service) UpdateByID(postID int, newPartialPost post.Post, fileBytes []byte) (post.Post, error) {
	oldPost, err := s.repository.GetByID(postID)
	if err != nil {
		err = fmt.Errorf("cannot update post with id %d: does not exist", postID)
		logger.Log.Debug(err.Error())
		return post.Post{}, err
	}

	dbPost, err := s.repository.UpdateByID(postID, newPartialPost)
	if err != nil {
		return post.Post{}, err
	}

	if len(fileBytes) > 0 {
		if err = s.bucket.Write(dbPost.Picture, fileBytes); err != nil {
			return post.Post{}, err
		}

		if err = s.bucket.Delete(oldPost.Picture); err != nil {
			return post.Post{}, err
		}
	}

	return dbPost, nil
}

// DeleteByID deletes a post by ID
func (s Service) DeleteByID(postID int) error {
	dbPost, err := s.repository.GetByID(postID)
	if err != nil {
		err = fmt.Errorf("cannot delete post with id %d: does not exist", postID)
		logger.Log.Debug(err.Error())
		return err
	}

	if err = s.bucket.Delete(dbPost.Picture); err != nil {
		return err
	}

	if err = s.repository.DeleteByID(postID); err != nil {
		return err
	}

	return nil
}
