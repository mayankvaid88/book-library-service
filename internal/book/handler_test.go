package book_test

import (
	"book-store/internal/book"
	mock_book "book-store/internal/mocks"
	"bytes"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type BookHandlerTestSuite struct {
	suite.Suite
	bookHandler book.BookHandler
	mockService *mock_book.MockBookService
	ctrl        *gomock.Controller
}

func TestBookHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BookHandlerTestSuite))
}

func (m *BookHandlerTestSuite) SetupSuite() {
	m.ctrl = gomock.NewController(m.Suite.T())
	m.mockService = mock_book.NewMockBookService(m.ctrl)
	m.bookHandler = *book.NewBookHandler(m.mockService)
}

func (m *BookHandlerTestSuite) TearDownTest() {
	m.ctrl.Finish()
}


func (m *BookHandlerTestSuite) TestList() {
        books := []book.Book{
            {ID: 1, Title: "A", Author: "X", Description: "desc A"},
            {ID: 2, Title: "B", Author: "Y", Description: "desc B"},
        }

    req,_ := http.NewRequest("GET", "/books?page=1&limit=5", nil)
    w := httptest.NewRecorder()

	m.mockService.EXPECT().List(req.Context(),5,0).Return(books,2,nil)

    m.bookHandler.List(w, req)

    res := w.Result()
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        m.Suite.T().Fatalf("expected 200, got %d", res.StatusCode)
    }

    var body struct {
        Page       int                   `json:"page"`
        Limit      int                   `json:"limit"`
        Total      int                   `json:"total"`
        TotalPages int                   `json:"totalPages"`
        Data       []book.BookResponse `json:"data"`
    }
    if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
        m.Suite.T().Fatalf("invalid JSON response: %v", err)
    }

    if body.Page != 1 || body.Limit != 5 {
        m.Suite.T().Errorf("bad page/limit: got %d/%d", body.Page, body.Limit)
    }
    if body.Total != 2 {
        m.Suite.T().Errorf("expected total=2, got %d", body.Total)
    }
    if body.TotalPages != int(math.Ceil(2/5.0)) {
        m.Suite.T().Errorf("bad totalPages: got %d", body.TotalPages)
    }
    if len(body.Data) != 2 {
        m.Suite.T().Errorf("expected 2 books, got %d", len(body.Data))
    }
}

func (m *BookHandlerTestSuite) TestList_ShouldSetLimitToThresholdWhenLimitPassedIsGreaterThanThreshold() {
        books := []book.Book{
            {ID: 1, Title: "A", Author: "X", Description: "desc A"},
            {ID: 2, Title: "B", Author: "Y", Description: "desc B"},
        }

    req,_ := http.NewRequest("GET", "/books?page=1&limit=500", nil)
    w := httptest.NewRecorder()

	m.mockService.EXPECT().List(req.Context(),100,0).Return(books,2,nil)

    m.bookHandler.List(w, req)

    res := w.Result()
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        m.Suite.T().Fatalf("expected 200, got %d", res.StatusCode)
    }

    var body struct {
        Page       int                   `json:"page"`
        Limit      int                   `json:"limit"`
        Total      int                   `json:"total"`
        TotalPages int                   `json:"totalPages"`
        Data       []book.BookResponse `json:"data"`
    }
    if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
        m.Suite.T().Fatalf("invalid JSON response: %v", err)
    }

    if body.Page != 1 || body.Limit != 100 {
        m.Suite.T().Errorf("bad page/limit: got %d/%d", body.Page, body.Limit)
    }
    if body.Total != 2 {
        m.Suite.T().Errorf("expected total=2, got %d", body.Total)
    }
    if body.TotalPages != int(math.Ceil(2/5.0)) {
        m.Suite.T().Errorf("bad totalPages: got %d", body.TotalPages)
    }
    if len(body.Data) != 2 {
        m.Suite.T().Errorf("expected 2 books, got %d", len(body.Data))
    }
}

func (m *BookHandlerTestSuite) TestList_ShouldThrowErrorWhenServiceReturnsError() {
    internalServerErr := book.GetErrorResponseByCode(book.InternalServerError)
    req := httptest.NewRequest("GET", "/books?page=1&limit=5", nil)
	m.mockService.EXPECT().List(req.Context(),5,0).Return(nil,0,internalServerErr)
	w := httptest.NewRecorder()
    m.bookHandler.List(w, req)

	bodyBytes, err := io.ReadAll(w.Result().Body)
	m.Suite.Nil(err)
	var actualErr book.ErrorResponse
	err = json.Unmarshal(bodyBytes, &actualErr)
	m.Suite.Nil(err)
	m.Suite.Equal(internalServerErr.Error(), actualErr.Error())
}

func (m *BookHandlerTestSuite) TestCreate() {
	createBookRequest := book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}
	requestBytes, _ := json.Marshal(createBookRequest)
	r, _ := http.NewRequest("GET", "/test/abcd", bytes.NewReader(requestBytes))
	w := httptest.NewRecorder()

	m.mockService.EXPECT().Create(r.Context(), createBookRequest).Return(int64(12), nil)

	m.bookHandler.Create(w, r)
	m.Suite.Equal(201, w.Result().StatusCode)
	m.Suite.Equal("/books/12", w.Result().Header.Get("Location"))
}

func (m *BookHandlerTestSuite) TestCreate_ShouldReturnErrorWhenServiceReturnsError() {
	createBookRequest := book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}
	requestBytes, _ := json.Marshal(createBookRequest)
	internalServerErr := book.GetErrorResponseByCode(book.InternalServerError)

	r, _ := http.NewRequest("GET", "/test/abcd", bytes.NewReader(requestBytes))
	w := httptest.NewRecorder()

	m.mockService.EXPECT().Create(r.Context(), createBookRequest).Return(int64(0), internalServerErr)

	m.bookHandler.Create(w, r)
	m.Suite.Equal(500, w.Result().StatusCode)

	bodyBytes, err := io.ReadAll(w.Result().Body)
	m.Suite.Nil(err)
	var actualErr book.ErrorResponse
	err = json.Unmarshal(bodyBytes, &actualErr)
	m.Suite.Nil(err)
	m.Suite.Equal(internalServerErr.Error(), actualErr.Error())
}

func (m *BookHandlerTestSuite) TestCreate_ShouldReturnBadRequestWhenRequestIsInvalid() {
	createBookRequest := book.CreateOrUpdateBookRequest{
		Title:       "",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}
	requestBytes, _ := json.Marshal(createBookRequest)

	r, _ := http.NewRequest("GET", "/test/abcd", bytes.NewReader(requestBytes))
	w := httptest.NewRecorder() 
	m.bookHandler.Create(w, r)
	m.Suite.Equal(400, w.Result().StatusCode)

	bodyBytes, err := io.ReadAll(w.Result().Body)
	m.Suite.Nil(err)
	var actualErr book.ErrorResponse
	err = json.Unmarshal(bodyBytes, &actualErr)
	m.Suite.Nil(err)
	m.Suite.Equal("Title failed on 'required'", actualErr.Error())
}

func (m *BookHandlerTestSuite) TestCreate_ShouldReturnErrorWhenRequestIsInvalid() {

	badReqErr := book.GetErrorResponseByCode(book.BadRequest)

	r, _ := http.NewRequest("GET", "/test/abcd", bytes.NewBuffer([]byte("invalid json")))
	w := httptest.NewRecorder()

	m.bookHandler.Create(w, r)
	m.Suite.Equal(400, w.Result().StatusCode)

	bodyBytes, err := io.ReadAll(w.Result().Body)
	m.Suite.Nil(err)
	var actualErr book.ErrorResponse
	err = json.Unmarshal(bodyBytes, &actualErr)
	m.Suite.Nil(err)
	m.Suite.Equal(badReqErr.Error(), actualErr.Error())
}

func (m *BookHandlerTestSuite) TestUpdate() {
	b := book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}

	responseBytes, _ := json.Marshal(b)

	r, _ := http.NewRequest("GET", "/books", bytes.NewBuffer(responseBytes))
	r = mux.SetURLVars(r, map[string]string{"id": "12"})
	m.mockService.EXPECT().CreateOrUpdate(r.Context(), 12, b).Return(int64(0), nil)
	w := httptest.NewRecorder()

	m.bookHandler.Update(w, r)
	m.Suite.Equal(204, w.Result().StatusCode)
}

func (m *BookHandlerTestSuite) TestUpdate_ShouldAddLocationHeaderWhenNewBookIsCreated() {
	b := book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}

	responseBytes, _ := json.Marshal(b)

	r, _ := http.NewRequest("GET", "/books", bytes.NewBuffer(responseBytes))
	r = mux.SetURLVars(r, map[string]string{"id": "12"})
	m.mockService.EXPECT().CreateOrUpdate(r.Context(), 12, b).Return(int64(12), nil)
	w := httptest.NewRecorder()

	m.bookHandler.Update(w, r)
	m.Suite.Equal(204, w.Result().StatusCode)
	m.Suite.Equal("/books/12", w.Result().Header.Get("Location"))
}

func (m *BookHandlerTestSuite) TestUpdate_ShouldThrowErrorWhenBookIdIsInvalid() {
	b := book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}

	responseBytes, _ := json.Marshal(b)

	r, _ := http.NewRequest("GET", "/books", bytes.NewBuffer(responseBytes))
	r = mux.SetURLVars(r, map[string]string{"id": "abcd"})
	w := httptest.NewRecorder()

	m.bookHandler.Update(w, r)
	m.Suite.Equal(400, w.Result().StatusCode)
}

func (m *BookHandlerTestSuite) TestUpdate_ShouldThrowErrorWhenServiceReturnsError() {
	b := book.CreateOrUpdateBookRequest{
		Title:       "Harry Potter",
		Author:      "JK Rolling",
		Description: "HarryPotter and Chambers of Secret",
	}

	requestBody, _ := json.Marshal(b)

	r, _ := http.NewRequest("PUT", "/books", bytes.NewBuffer(requestBody))
	r = mux.SetURLVars(r, map[string]string{"id": "12"})
	w := httptest.NewRecorder()
	internalErr := book.GetErrorResponseByCode(book.InternalServerError)
	m.mockService.EXPECT().CreateOrUpdate(r.Context(), 12, b).Return(int64(0),internalErr)

	m.bookHandler.Update(w, r)
	m.Suite.Equal(500, w.Result().StatusCode)

	var actualErr book.ErrorResponse
	bodyBytes, err := io.ReadAll(w.Result().Body)
	m.Suite.Nil(err)
	err = json.Unmarshal(bodyBytes, &actualErr)
	m.Suite.Nil(err)
	m.Suite.Equal(internalErr.Error(), actualErr.Error())
}

func (m *BookHandlerTestSuite) TestDelete() {
	r, _ := http.NewRequest("DELETE", "/books", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "12"})
	w := httptest.NewRecorder()
	m.mockService.EXPECT().Delete(r.Context(), 12).Return(nil)
	m.bookHandler.Delete(w, r)
	m.Suite.Equal(204, w.Result().StatusCode)
}

func (m *BookHandlerTestSuite) TestDelete_ShouldReturnErrorWhenServiceReturnsError() {
	r, _ := http.NewRequest("DELETE", "/books", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "12"})
	w := httptest.NewRecorder()
	internalErr := book.GetErrorResponseByCode(book.InternalServerError)
	m.mockService.EXPECT().Delete(r.Context(), 12).Return(internalErr)
	m.bookHandler.Delete(w, r)
	var actualErr book.ErrorResponse
	bodyBytes, err := io.ReadAll(w.Result().Body)
	m.Suite.Nil(err)
	err = json.Unmarshal(bodyBytes, &actualErr)
	m.Suite.Nil(err)
	m.Suite.Equal(internalErr.Error(), actualErr.Error())
}
