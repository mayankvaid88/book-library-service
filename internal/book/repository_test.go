package book_test

import (
	"book-store/internal/book"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type BookRepositoryTestSuite struct {
	suite.Suite
	bookRepository book.BookRepository
	sqlMock        sqlmock.Sqlmock
	db             *sql.DB
	ctrl           *gomock.Controller
}

func TestBookRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BookRepositoryTestSuite))
}

func (m *BookRepositoryTestSuite) SetupSuite() {
	m.ctrl = gomock.NewController(m.T())
	m.db, m.sqlMock, _ = sqlmock.New()
	m.bookRepository = book.NewBookRepository(m.db)
}

func (m *BookRepositoryTestSuite) TearDownTest() {
	m.ctrl.Finish()
}

func (m *BookRepositoryTestSuite) TestCreate_ShouldInsertBookInDatabase() {
	b := book.Book{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}
	m.sqlMock.ExpectQuery(regexp.QuoteMeta("INSERT INTO books (title, author, description) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs("Harry Potter", "JK Rolling", "HarryPotter and Chambers of Secret").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(10))
	bId, err := m.bookRepository.Create(context.Background(), b)
	m.Suite.Nil(m.sqlMock.ExpectationsWereMet())
	m.Suite.Nil(err)
	m.Suite.Equal(int64(10), bId)
}

func (m *BookRepositoryTestSuite) TestCreate_ShouldThrowErrorWhenInsertQueryFails() {
	b := book.Book{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}
	m.sqlMock.ExpectQuery(regexp.QuoteMeta("INSERT INTO books (title, author, description) VALUES ($1, $2, $3) RETURNING id")).
		WithArgs("Harry Potter", "JK Rolling", "HarryPotter and Chambers of Secret").
		WillReturnError(errors.New("unique constraint violation"))
	bId, err := m.bookRepository.Create(context.Background(), b)
	m.Suite.Equal(int64(0), bId)
	m.Suite.Nil(m.sqlMock.ExpectationsWereMet())
	m.Suite.EqualError(err, "unique constraint violation")
}

func (m *BookRepositoryTestSuite) TestGetById_ShouldShouldReturnBookWithTheProvidedId() {
	rows := sqlmock.NewRows([]string{"id", "title", "author", "description"}).
		AddRow(12, "Harry Potter", "JK Rolling", "HarryPotter and Chambers of Secret")
	m.sqlMock.ExpectQuery("SELECT id, title, author, description FROM books").WillReturnRows(rows)
	b, err := m.bookRepository.GetByID(context.Background(), 12)
	m.Suite.Nil(m.sqlMock.ExpectationsWereMet())
	m.Suite.Nil(err)
	m.Suite.Equal(book.Book{
		ID:          12,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}, b)
}

func (m *BookRepositoryTestSuite) TestGetById_ShouldReturnNotFoundErrorIfNoBookPresentForGivenId() {
	m.sqlMock.ExpectQuery("SELECT id, title, author, description FROM books").WillReturnError(sql.ErrNoRows)
	b, err := m.bookRepository.GetByID(context.Background(), 12)
	m.Suite.Nil(m.sqlMock.ExpectationsWereMet())
	m.Suite.Empty(b)
	m.Suite.EqualError(err, "book not found")
}

func (m *BookRepositoryTestSuite) TestList_ShouldReturAllBooks() {
	rows := sqlmock.NewRows([]string{"id", "title", "author", "description","total_count"}).
		AddRow(12, "Harry Potter", "JK Rolling", "HarryPotter and Chambers of Secret",2).
		AddRow(13, "Harry Potter", "JK Rolling", "HarryPotter and Goblet of Fire",2)
	m.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, description, COUNT(*) OVER() AS total_count FROM books ORDER BY id LIMIT $1 OFFSET $2")).WithArgs(5,1).WillReturnRows(rows)
	b,totalCount, err := m.bookRepository.List(context.Background(),5,1)
		m.Suite.Nil(err)
	m.Suite.Nil(m.sqlMock.ExpectationsWereMet())
	m.Suite.Len(b, 2)
	m.Suite.Equal(2,totalCount)
	m.Suite.Equal([]book.Book{{
		ID:          12,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}, {
		ID:          13,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Goblet of Fire",
	}}, b)
}

func (m *BookRepositoryTestSuite) TestList_ShouldReturnErrorWhenQueryFails() {
	m.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, author, description, COUNT(*) OVER() AS total_count FROM books ORDER BY id LIMIT $1 OFFSET $2")).WithArgs(1,2).WillReturnError(errors.New("unable to connect"))
	b,totalCount ,err := m.bookRepository.List(context.Background(),1,2)
	m.Suite.Nil(m.sqlMock.ExpectationsWereMet())
	m.Suite.Equal(0,totalCount)
	m.Suite.EqualError(err, "unable to connect")
	m.Suite.Empty(b)
}

func (m *BookRepositoryTestSuite) TestUpdate_ShouldUpdateTheBookRecord() {
	m.sqlMock.ExpectExec("UPDATE books").
		WithArgs("Harry Potter", "JK Rolling", "HarryPotter and Goblet of Fire", 13).
		WillReturnResult(sqlmock.NewResult(13, 1))
	err := m.bookRepository.Update(context.Background(), book.Book{
		ID:          13,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Goblet of Fire",
	})
	m.Suite.Nil(m.sqlMock.ExpectationsWereMet())
	m.Suite.Nil(err)
}

func (m *BookRepositoryTestSuite) TestUpdate_ShouldReturnErrorWhenUpdateFails() {
	m.sqlMock.ExpectExec("UPDATE books").WillReturnError(errors.New("unable to connect"))
	err := m.bookRepository.Update(context.Background(), book.Book{
		ID:          13,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Goblet of Fire",
	})
	m.Suite.Nil(m.sqlMock.ExpectationsWereMet())
	m.Suite.EqualError(err, "unable to connect")
}
