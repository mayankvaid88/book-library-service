package book

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("book not found")

type BookRepository interface {
	Create(ctx context.Context, b Book) (int64, error)
	GetByID(ctx context.Context, id int) (Book, error)
	List(ctx context.Context,limit, offset int) ([]Book,int, error)
	Update(ctx context.Context, b Book) error
	Delete(ctx context.Context, id int) error
}

type sqlBookRepo struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &sqlBookRepo{db: db}
}

func (r *sqlBookRepo) Create(ctx context.Context, b Book) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO books (title, author, description) VALUES ($1, $2, $3) RETURNING id`,
		b.Title, b.Author, b.Description).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *sqlBookRepo) GetByID(ctx context.Context, id int) (Book, error) {
	b := Book{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, author, description FROM books WHERE id = $1`, id).
		Scan(&b.ID, &b.Title, &b.Author, &b.Description)
	if err == sql.ErrNoRows {
		return Book{}, ErrNotFound
	}
	return b, err
}

func (r *sqlBookRepo) List(ctx context.Context,limit, offset int) ([]Book,int, error) {
	  rows, err := r.db.QueryContext(ctx, `
        SELECT id, title, author, description,
               COUNT(*) OVER() AS total_count
        FROM books
        ORDER BY id
        LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil,0, err
	}
	defer rows.Close()
	var books []Book
	var total int
	for rows.Next() {
		b := Book{}
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Description,&total); err != nil {
			return nil,0, err
		}
		books = append(books, b)
	}
	return books,total, nil
}

func (r *sqlBookRepo) Update(ctx context.Context, b Book) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE books SET title=$1, author=$2, description=$3 WHERE id=$4`,
		b.Title, b.Author, b.Description, b.ID)
	return err
}

func (r *sqlBookRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM books WHERE id=$1`, id)
	return err
}
