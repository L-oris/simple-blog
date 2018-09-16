package postrepository

import (
	"database/sql"

	"github.com/L-oris/yabb/logger"
	"github.com/L-oris/yabb/models/db"
	"github.com/L-oris/yabb/models/post"
)

type Repository struct {
	DB *sql.DB
}

func New() *Repository {
	return &Repository{
		DB: db.BlogDB,
	}
}

func (r Repository) Ping() error {
	if err := r.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (r Repository) GetAll() ([]post.Post, error) {
	rows, err := r.DB.Query("SELECT * FROM Posts;")
	if err != nil {
		logger.Log.Error("query error: ", err)
		return make([]post.Post, 0), err
	}
	defer rows.Close()

	var result []post.Post
	for rows.Next() {
		post := post.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			logger.Log.Error("scan error: ", err)
			return nil, err
		}
		result = append(result, post)
	}
	if err := rows.Err(); err != nil {
		logger.Log.Error("rows error: ", err)
		return nil, err
	}

	return result, nil
}
