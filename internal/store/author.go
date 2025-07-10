package store

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Author struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Bio       string    `db:"bio"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type AuthorRepo struct {
	db *sqlx.DB
}

func NewAuthorRepo(db *sqlx.DB) *AuthorRepo {
	return &AuthorRepo{db: db}
}

func (r *AuthorRepo) CreateAuthor(author *Author) error {
	q := `
INSERT INTO 
authors 
	(name, email, bio, created_at, updated_at) 
VALUES 
	($1, $2, $3, $4, $5)
`

	res, err := r.db.Exec(q,
		author.Name,
		author.Email,
		author.Bio,
		author.CreatedAt.Format(time.RFC3339),
		author.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	author.ID = id

	return nil
}

func (r *AuthorRepo) ReadAllAuthors() ([]*Author, error) {
	var authors []*Author

	q := "SELECT * FROM authors ORDER BY created_at DESC"

	err := r.db.Select(&authors, q)
	if err != nil {
		return nil, fmt.Errorf("Error getting all authors: %w", err)
	}

	return authors, nil
}

func (r *AuthorRepo) ReadAuthorByID(id int64) (*Author, error) {
	var author Author
	err := r.db.Get(&author, "SELECT * FROM authors WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("Error getting author by id `%d`: %w", id, err)
	}

	return &author, nil
}

func (r *AuthorRepo) UpdateAuthor(author *Author) error {
	q := `
UPDATE 
	authors 
SET 
	name = ?, 
	email = ?, 
	bio = ?, 
	updated_at = ?
WHERE 
	id = ?
`
	_, err := r.db.Exec(q,
		author.Name,
		author.Email,
		author.Bio,
		author.UpdatedAt,
		author.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthorRepo) DeleteAuthorByID(id int64) error {
	_, err := r.db.Exec("DELETE FROM authors WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}