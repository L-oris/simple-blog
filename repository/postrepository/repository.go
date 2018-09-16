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

// New creates a new Repository
func New() *Repository {
	return &Repository{
		DB: db.BlogDB,
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
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			logger.Log.Error("scan error: ", err)
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
	if err := row.Scan(&result.ID, &result.Title, &result.Content, &result.CreatedAt); err != nil {
		logger.Log.Warning("scan error: ", err)
		return post.Post{}, err
	}

	return result, nil
}

// Add adds new Post to DB
// The following fields cannot be managed externally: ID, CreatedAt
// Returns the new Post
func (r Repository) Add(partialPost post.Post) (post.Post, error) {
	sqlStatement := `INSERT INTO Posts (Title, Content) VALUES (?, ?);`
	queryReturn, err := r.DB.Exec(sqlStatement, partialPost.Title, partialPost.Content)
	if err != nil {
		logger.Log.Warning("insert error: ", err)
		return post.Post{}, err
	}

	lastInsertID, _ := queryReturn.LastInsertId()
	return r.GetByID(int(lastInsertID))
}