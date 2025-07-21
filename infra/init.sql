\c book_store
CREATE TABLE books (
  id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title       VARCHAR(255) NOT NULL,
  author      VARCHAR(255) NOT NULL,
  description TEXT NOT NULL
);