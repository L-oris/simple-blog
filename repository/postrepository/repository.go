package postrepository

import (
	"database/sql"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/post"
	"github.com/imdario/mergo"
)

type Config struct {
	DB *sql.DB
}

type Repository struct {
	DB *sql.DB
}

// New creates a new Repository
func New(config *Config) *Repository {
	return &Repository{
		DB: config.DB,
	}
}

// Ping checks DB connection
func (r Repository) Ping() error {
	if err := r.DB.Ping(); err != nil {
		return err
	}
	return nil
}

// GetAll gets all Posts
func (r Repository) GetAll() ([]post.Post, error) {
	sqlStatement := `SELECT * FROM Posts;`
	rows, _ := r.DB.Query(sqlStatement)
	defer rows.Close()

	var result []post.Post
	for rows.Next() {
		post := post.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ImageID, &post.CreatedAt); err != nil {
			logger.Log.Error("scan error: %s", err.Error())
			return nil, err
		}
		result = append(result, post)
	}

	return result, nil
}

// GetByID gets Post by ID
// Returns error when not found
func (r Repository) GetByID(id int) (post.Post, error) {
	sqlStatement := `SELECT * FROM Posts WHERE ID=?;`
	row := r.DB.QueryRow(sqlStatement, id)

	result := post.Post{}
	if err := row.Scan(&result.ID, &result.Title, &result.Content, &result.ImageID, &result.CreatedAt); err != nil {
		logger.Log.Warning("scan error: %s", err.Error())
		return post.Post{}, err
	}

	return result, nil
}

// Add adds new Post to DB
// The following fields cannot be managed externally: ID, CreatedAt
// Returns the new Post
func (r Repository) Add(partialPost post.Post) (post.Post, error) {
	sqlStatement := `INSERT INTO Posts (Title, Content, ImageID) VALUES (?, ?, ?);`
	queryReturn, err := r.DB.Exec(sqlStatement, partialPost.Title, partialPost.Content, partialPost.ImageID)
	if err != nil {
		logger.Log.Warning("insert error: %s", err.Error())
		return post.Post{}, err
	}

	lastInsertID, _ := queryReturn.LastInsertId()
	return r.GetByID(int(lastInsertID))
}

// UpdateByID updates only provided fields in existing DB Post
func (r Repository) UpdateByID(id int, partialPost post.Post) (post.Post, error) {
	dbPost, err := r.GetByID(id)
	if err != nil {
		logger.Log.Warning("post ", string(id), "not found")
		return post.Post{}, nil
	}

	if err = mergo.Merge(&dbPost, partialPost, mergo.WithOverride); err != nil {
		logger.Log.Error("failed to merge posts: %s", err.Error())
		return post.Post{}, err
	}

	sqlStatement := `UPDATE Posts SET Title=?, Content=?, ImageID=? WHERE ID=?`
	_, err = r.DB.Exec(sqlStatement, dbPost.Title, dbPost.Content, dbPost.ImageID, dbPost.ID)
	if err != nil {
		logger.Log.Error("cannot update post: %s", err.Error())
		return post.Post{}, nil
	}

	return r.GetByID(id)
}

// DeleteByID deletes post by ID
// Returns an error if Post is not found
func (r Repository) DeleteByID(id int) error {
	_, err := r.GetByID(id)
	if err != nil {
		logger.Log.Warning("post ", string(id), "not found")
		return err
	}

	sqlStatement := `DELETE FROM Posts WHERE ID=?`
	_, err = r.DB.Exec(sqlStatement, id)
	if err != nil {
		logger.Log.Error("cannot delete post: %s", err.Error())
		return err
	}

	return nil
}
