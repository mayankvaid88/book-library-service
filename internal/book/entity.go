package book

type Book struct {
	ID          int    `sql:"id"`
	Title       string `sql:"title"`
	Author      string `sql:"author"`
	Description string `sql:"description"`
}
