package book

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type BookHandler struct {
	svc BookService
	val validator.Validate
}

func NewBookHandler(s BookService) *BookHandler {
	return &BookHandler{svc: s,val: *validator.New()}
}

// List godoc
// @Summary      List books with pagination
// @Description  Returns a paginated list of books
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        page   query     int  false  "Page number (default 1)"    default(1)
// @Param        limit  query     int  false  "Page size (1–100, default 10)" default(10)
// @Success      200    {object}  PaginatedBookListResponse
// @Failure      400    {object}  ErrorResponse
// @Router       /books [get]
func (h *BookHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, pageConvErr := strconv.Atoi(q.Get("page"))
	if pageConvErr!=nil{
		logrus.Error("invalid page number provided ",q.Get("page"))
		sendError(w, *GetErrorResponseByCode(BadRequest))
		return
	}
	limit, limitConvErr := strconv.Atoi(q.Get("limit"))
	if limitConvErr!=nil{
		logrus.Error("invalid limit number provided ",q.Get("limit"))
		sendError(w,*GetErrorResponseByCode(BadRequest))
	}
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	offset := (page - 1) * limit

	books,totalCount, err := h.svc.List(r.Context(),limit,offset)
	if err != nil {
		sendError(w, *err)
		return
	}
	out := make([]BookResponse, len(books))
	for i, b := range books {
		out[i] = BookResponse{ID: b.ID, Title: b.Title, Author: b.Author, Description: b.Description}
	}
	p := PaginatedBookListResponse{
		Page: page,
		Limit: limit,
		Total: totalCount,
		TotalPages: int(math.Ceil(float64(totalCount) / float64(limit))),
		Data: out,
	}
	json.NewEncoder(w).Encode(p)
}


// Get godoc
// @Summary      Get book by ID
// @Description  Retrieve a single book by its ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id     path      int   true   "Book ID"
// @Success      200    {object}  BookResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      404    {object}  ErrorResponse
// @Router       /books/{id} [get]
func (h *BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.Atoi(mux.Vars(r)["id"])
	if convErr!=nil{
		logrus.Error("invalid book id provided ",mux.Vars(r)["id"])
		sendError(w, *GetErrorResponseByCode(BadRequest))
		return
	}
	b, err := h.svc.Get(r.Context(), id)
	if err != nil {
		sendError(w, *err)
		return
	}
	json.NewEncoder(w).Encode(BookResponse{ID: b.ID, Title: b.Title, Author: b.Author, Description: b.Description})
}


// Create godoc
// @Summary      Create a new book
// @Description  Add a new book record
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      CreateOrUpdateBookRequest  true  "Book data"
// @Success      201    {object}  nil
// @Header       201    {string}  Location  "URL of created book"
// @Failure      400    {object}  ErrorResponse
// @Router       /books [post]
func (h *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateOrUpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, *GetErrorResponseByCode(BadRequest))
		return
	}

    if err := h.val.Struct(&req); err != nil {
        var errs []string
        for _, fe := range err.(validator.ValidationErrors) {
            errs = append(errs, fmt.Sprintf("%s failed on '%s'", fe.Field(), fe.Tag()))
        }
		logrus.Error("error while validating the request. error is ",errs)
        sendError(w, *GetErrorResponse(BadRequest, strings.Join(errs, "; "), http.StatusBadRequest))
        return
    }

	bId, err := h.svc.Create(r.Context(), req)
	if err != nil {
		sendError(w, *err)
		return
	}
	w.Header().Set("location", fmt.Sprintf("%s/%d", "/books", bId))
	w.WriteHeader(http.StatusCreated)
}

// Update godoc
// @Summary      Update or create book by ID
// @Description  Update existing book or create if not exists
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int                        true  "Book ID"
// @Param        book  body      CreateOrUpdateBookRequest  true  "Book data"
// @Success      204    {object}  nil
// @Header       204    {string}  Location  "Optional new resource URL"
// @Failure      400    {object}  ErrorResponse
// @Failure      404    {object}  ErrorResponse
// @Router       /books/{id} [put]
func (h *BookHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req CreateOrUpdateBookRequest
	id, cErr := strconv.Atoi(mux.Vars(r)["id"])
	if cErr != nil {
		logrus.Error("invalid book id provided ",mux.Vars(r)["id"])
		sendError(w, *GetErrorResponseByCode(BadRequest))
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, *GetErrorResponseByCode(BadRequest))
		return
	}
	bId, err := h.svc.CreateOrUpdate(r.Context(), id, req)
	if err != nil {
		sendError(w, *err)
		return
	}
	if bId != 0 {
		w.Header().Set("location", fmt.Sprintf("%s/%d", "/books", bId))
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary      Delete book by ID
// @Description  Remove a book record
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int   true   "Book ID"
// @Success      204    {object}  nil
// @Failure      400    {object}  ErrorResponse
// @Failure      404    {object}  ErrorResponse
// @Router       /books/{id} [delete]
func (h *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, convErr := strconv.Atoi(mux.Vars(r)["id"])
	if convErr!=nil{
		logrus.Error("invalid book id provided ",mux.Vars(r)["id"])
		sendError(w, *GetErrorResponseByCode(BadRequest))
		return
	}
	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		sendError(w, *err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func sendError(w http.ResponseWriter, errResponse ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResponse.HttpStatusCode)
	json.NewEncoder(w).Encode(errResponse)
}
