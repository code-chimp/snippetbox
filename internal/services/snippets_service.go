package services

import (
	"database/sql"
	"errors"
	"github.com/code-chimp/snippetbox/internal/models"
)

// SnippetsService wraps a sql.DB connection pool.
type SnippetsService struct {
	DB *sql.DB
}

// Insert will insert a new snippet into the database.
func (m *SnippetsService) Insert(title string, content string, expires int) (int, error) {
	query := `INSERT INTO snippets
		      (title, content, created, expires)
      	      VALUES
      	      (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get will return a specific snippet based on its id.
func (m *SnippetsService) Get(id int) (models.Snippet, error) {
	query := `SELECT id, title, content, created, expires
              FROM snippets
              WHERE expires > UTC_TIMESTAMP() AND id = ?`

	var s models.Snippet

	err := m.DB.QueryRow(query, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Snippet{}, models.ErrNoRecord
		} else {
			return models.Snippet{}, err
		}
	}

	return s, nil
}

// Latest will return the 10 most recently created snippets.
func (m *SnippetsService) Latest() ([]models.Snippet, error) {
	query := `SELECT id, title, content, created, expires
              FROM snippets
              WHERE expires > UTC_TIMESTAMP()
              ORDER BY id DESC
              LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []models.Snippet

	for rows.Next() {
		var s models.Snippet

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
