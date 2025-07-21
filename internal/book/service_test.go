package book_test

import (
	"book-store/internal/book"
	mock_book "book-store/internal/mocks"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type BookServiceTestSuite struct {
	suite.Suite
	bookService book.BookService
	mockRepo    *mock_book.MockBookRepository
	ctrl        *gomock.Controller
}

func TestBookServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BookServiceTestSuite))
}

func (m *BookServiceTestSuite) SetupSuite() {
	m.ctrl = gomock.NewController(m.Suite.T())
	m.mockRepo = mock_book.NewMockBookRepository(m.ctrl)
	m.bookService = book.NewBookService(m.mockRepo)
}

func (m *BookServiceTestSuite) TearDownTest() {
	m.ctrl.Finish()
}

func (m *BookServiceTestSuite) TestCreate() {
	m.mockRepo.EXPECT().Create(context.Background(), book.Book{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}).Return(int64(12), nil)
	bId, err := m.bookService.Create(context.Background(), book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	})
	m.Suite.Nil(err)
	m.Suite.Equal(bId, int64(12))
}

func (m *BookServiceTestSuite) TestCreate_ShouldReturnErrorWhenRepositoryFails() {
	m.mockRepo.EXPECT().Create(context.Background(), book.Book{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}).Return(int64(0), errors.New("unable to connect"))
	_, err := m.bookService.Create(context.Background(), book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	})
	m.Suite.Equal(err, book.GetErrorResponseByCode(book.InternalServerError))
}

func (m *BookServiceTestSuite) TestGet_ShouldReturnBookForGivenId() {
	m.mockRepo.EXPECT().GetByID(context.Background(), 12).Return(book.Book{
		ID:          12,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}, nil)
	b, err := m.bookService.Get(context.Background(), 12)
	m.Suite.Nil(err)
	m.Suite.Equal(b, book.Book{
		ID:          12,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	})
}

func (m *BookServiceTestSuite) TestGet_ShouldReturnNotFoundIfBookWithGivenIDDoesNotExist() {
	m.mockRepo.EXPECT().GetByID(context.Background(), 12).Return(book.Book{}, book.ErrNotFound)
	b, err := m.bookService.Get(context.Background(), 12)
	m.Suite.Equal(err, book.GetErrorResponseByCode(book.BookNotFound))
	m.Suite.Empty(b)
}

func (m *BookServiceTestSuite) TestGet_ShouldReturnInternalServerErrorIfRepositoryFails() {
	m.mockRepo.EXPECT().GetByID(context.Background(), 12).Return(book.Book{}, errors.New("unable to connect"))
	b, err := m.bookService.Get(context.Background(), 12)
	m.Suite.Equal(err, book.GetErrorResponseByCode(book.InternalServerError))
	m.Suite.Empty(b)
}

func (m *BookServiceTestSuite) TestList_ShouldReturnAllBooksForCurrentPage() {
	m.mockRepo.EXPECT().List(context.Background(), 10, 2).Return([]book.Book{{
		ID:          12,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}, {
		ID:          13,
		Title:       "Harry Potter 2",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret 2",
	},
	},1, nil)
	b, totalCount, err := m.bookService.List(context.Background(), 10, 2)
	m.Suite.Nil(err)
	m.Suite.Equal(1,totalCount)
	m.Suite.Equal(b, []book.Book{{
		ID:          12,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}, {
		ID:          13,
		Title:       "Harry Potter 2",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret 2",
	},
	})
}

func (m *BookServiceTestSuite) TestList_ShouldReturnErrorWhenRepositoryFails() {
	m.mockRepo.EXPECT().List(context.Background(), 10, 2).Return(nil,0, errors.New("unable to connect"))
	b,totalCount, err := m.bookService.List(context.Background(), 10, 2)
	m.Suite.Nil(b)
	m.Suite.Zero(totalCount)
	m.Suite.Equal(book.GetErrorResponseByCode(book.InternalServerError),err)
}

func (m *BookServiceTestSuite) TestUpdate() {
	b := book.Book{
		ID:          12,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}
	m.mockRepo.EXPECT().GetByID(context.Background(), 12).Return(b, nil)
	m.mockRepo.EXPECT().Update(context.Background(), book.Book{
		ID:          12,
		Title:       "Harry Potter 4",
		Author:      "JKR",
		Description: "HarryPotter and Goblet Of Fire",
	}).Return(nil)
	bId, err := m.bookService.CreateOrUpdate(context.Background(), 12, book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter 4",
		Author:      "JKR",
		Description: "HarryPotter and Goblet Of Fire",
	})
	m.Suite.Zero(bId)
	m.Suite.Nil(err)
}

func (m *BookServiceTestSuite) TestUpdate_ShouldCreateNewBookWhenBookWithGivenIdDoesNotExit() {
	m.mockRepo.EXPECT().GetByID(context.Background(), 12).Return(book.Book{}, book.ErrNotFound)
	m.mockRepo.EXPECT().Create(context.Background(), book.Book{
		Title:       "Harry Potter 4",
		Author:      "JKR",
		Description: "HarryPotter and Goblet Of Fire",
	}).Return(int64(1), nil)
	bId, err := m.bookService.CreateOrUpdate(context.Background(), 12, book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter 4",
		Author:      "JKR",
		Description: "HarryPotter and Goblet Of Fire",
	})
	m.Suite.Equal(int64(1), bId)
	m.Suite.Nil(err)
}

func (m *BookServiceTestSuite) TestUpdate_ShouldReturnInternalServerErrorWhenUpdateFails() {
	b := book.Book{
		ID:          12,
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}
	m.mockRepo.EXPECT().GetByID(context.Background(), 12).Return(b, nil)
	m.mockRepo.EXPECT().Update(context.Background(), book.Book{
		ID:          12,
		Title:       "Harry Potter 4",
		Author:      "JKR",
		Description: "HarryPotter and Goblet Of Fire",
	}).Return(errors.New("unable to connect"))
	bId, err := m.bookService.CreateOrUpdate(context.Background(), 12, book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter 4",
		Author:      "JKR",
		Description: "HarryPotter and Goblet Of Fire",
	})
	m.Suite.Zero(bId)
	m.Suite.Equal(err, book.GetErrorResponseByCode(book.InternalServerError))
}

func (m *BookServiceTestSuite) TestDelete() {
	m.mockRepo.EXPECT().Delete(context.Background(), 12).Return(nil)
	err := m.bookService.Delete(context.Background(), 12)
	m.Suite.Nil(err)
}

func (m *BookServiceTestSuite) TestDelete_ShouldReturnInternalServerErrorWhenUpdateFails() {
	m.mockRepo.EXPECT().Delete(context.Background(), 12).Return(errors.New("unable to connect"))
	err := m.bookService.Delete(context.Background(), 12)
	m.Suite.Equal(err, book.GetErrorResponseByCode(book.InternalServerError))
}
