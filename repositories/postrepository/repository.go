package postrepository

import (
	"database/sql"
	"errors"
	"fmt"

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
func New(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
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
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Picture, &post.CreatedAt); err != nil {
			logger.Log.Warningf("scan error: %s", err.Error())
			return nil, errors.New("cannot get posts")
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
	if err := row.Scan(&result.ID, &result.Title, &result.Content, &result.Picture, &result.CreatedAt); err != nil {
		logger.Log.Warning("scan error: %s", err.Error())
		return post.Post{}, fmt.Errorf("cannot get post %s", string(id))
	}

	return result, nil
}

// Add adds new Post to DB
// The following fields cannot be managed externally: ID, CreatedAt
// Returns the new Post
func (r Repository) Add(partialPost post.Post) (post.Post, error) {
	sqlStatement := `INSERT INTO Posts (Title, Content, Picture) VALUES (?, ?, ?);`
	queryReturn, err := r.DB.Exec(sqlStatement, partialPost.Title, partialPost.Content, partialPost.Picture)
	if err != nil {
		logger.Log.Warning("insert error: %s", err.Error())
		return post.Post{}, errors.New("cannot add new post")
	}

	lastInsertID, _ := queryReturn.LastInsertId()
	return r.GetByID(int(lastInsertID))
}

// UpdateByID updates only provided fields in existing DB Post
// Returns updated Post
func (r Repository) UpdateByID(id int, partialPost post.Post) (post.Post, error) {
	defaultError := fmt.Errorf("cannot update post %d", id)

	dbPost, err := r.GetByID(id)
	if err != nil {
		return post.Post{}, err
	}

	if err = mergo.Merge(&dbPost, partialPost, mergo.WithOverride); err != nil {
		logger.Log.Warningf("merge posts error: %s", err.Error())
		return post.Post{}, defaultError
	}

	sqlStatement := `UPDATE Posts SET Title=?, Content=?, Picture=? WHERE ID=?`
	_, err = r.DB.Exec(sqlStatement, dbPost.Title, dbPost.Content, dbPost.Picture, dbPost.ID)
	if err != nil {
		logger.Log.Error("update post error: %s", err.Error())
		return post.Post{}, defaultError
	}

	return r.GetByID(id)
}

// DeleteByID deletes post by ID
func (r Repository) DeleteByID(id int) error {
	if _, err := r.GetByID(id); err != nil {
		return err
	}

	sqlStatement := `DELETE FROM Posts WHERE ID=?`
	if _, err := r.DB.Exec(sqlStatement, id); err != nil {
		logger.Log.Error("delete post error: %s", err.Error())
		return fmt.Errorf("cannot delete post %d", id)
	}

	return nil
}
