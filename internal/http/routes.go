package http

import (
	"book-store/internal/book"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, db *sql.DB) {
	bookRepo := book.NewBookRepository(db)
	bookService := book.NewBookService(bookRepo)
	handler := book.NewBookHandler(bookService)

	r.HandleFunc("/books", handler.List).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", handler.Get).Methods(http.MethodGet)
	r.HandleFunc("/books", handler.Create).Methods(http.MethodPost)
	r.HandleFunc("/books/{id}", handler.Update).Methods(http.MethodPut)
	r.HandleFunc("/books/{id}", handler.Delete).Methods(http.MethodDelete)
}
