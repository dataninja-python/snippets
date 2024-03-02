package models

import (
	"database/sql"
	"errors"
	"time"
)

// What is a Snippet?
// This defines in terms of its components as a Go object.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Defines the specific db object used by Snippets
type SnippetModel struct {
	DB *sql.DB
}

// Now we add actual Snippet functionality
// Question: why don't we use an interface?
// Insert a snippet into database
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	// SQL statement to execute
	stmt := `INSERT INTO snippetcol (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LstInsertId() method on the result to get the ID of our newly inserted record in the snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// Confirm the id is the proper int64 type before returning
	return int(id), nil
}

// Return A specific snippet by its id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	// SQL statement
	stmt := `SELECT id, title, content, created, expires FROM snippetcol WHERE expires > UTC_TIMESTAMP() AND id=?`

	// Use the QueryRow() method on the connection pool to execute the SQL statement
	// to get the row
	row := m.DB.QueryRow(stmt, id)

	// A new empty Snippet object
	var s Snippet

	// pull the desired data from the row
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// Return 10 most recent snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {

	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippetcol WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use Query() method on the connection pool to execute our SQL statement
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Defer closing the database
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet struct
	var snippets []Snippet

	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct
		var s Snippet

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Append i to the slice of snippets
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Everything worked
	return snippets, nil

}
