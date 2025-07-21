package integrationtest

import (
	"context"
	"testing"
)

func cleanUp(t *testing.T) {
	_, err := sharedDB.Exec(`TRUNCATE books RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("cleanup error: %v", err)
	}
}

func insertTestBook(ctx context.Context, title, author, description string) (int64, error) {
	const query = `
      INSERT INTO books (title, author, description)
      VALUES ($1, $2, $3)
      RETURNING id
    `
	var id int64
	err := sharedDB.
		QueryRowContext(ctx, query, title, author, description).
		Scan(&id)
	return id, err
}
