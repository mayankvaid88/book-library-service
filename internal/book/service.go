package book

import (
	"context"
	"github.com/sirupsen/logrus"
)

type BookService interface {
	Create(ctx context.Context, req CreateOrUpdateBookRequest) (int64, *ErrorResponse)
	Get(ctx context.Context, id int) (Book, *ErrorResponse)
	List(ctx context.Context,limit, offset int) ([]Book,int, *ErrorResponse)
	CreateOrUpdate(ctx context.Context, id int, req CreateOrUpdateBookRequest) (int64, *ErrorResponse)
	Delete(ctx context.Context, id int) *ErrorResponse
}

type bookService struct {
	repository BookRepository
}

func NewBookService(r BookRepository) BookService {
	return &bookService{repository: r}
}

func (s *bookService) Create(ctx context.Context, req CreateOrUpdateBookRequest) (int64, *ErrorResponse) {
	b := Book{
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
	}
	id, err := s.repository.Create(ctx, b)
	if err != nil {
		logrus.Error("error while creatin book. error is ",err)
		return 0, GetErrorResponseByCode(InternalServerError)
	}
	return id, nil
}

func (s *bookService) Get(ctx context.Context, id int) (Book, *ErrorResponse) {
	book, err := s.repository.GetByID(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			logrus.Error("no record found for given id ",id)
			return Book{}, GetErrorResponseByCode(BookNotFound)
		}
		logrus.Error("error while fetching the record for id ",id," error is ",err)
		return Book{}, GetErrorResponseByCode(InternalServerError)
	}
	return book, nil
}

func (s *bookService) List(ctx context.Context,limit, offset int) ([]Book,int, *ErrorResponse) {
	books,totalCount, err := s.repository.List(ctx, limit, offset)
	if err != nil {
		logrus.Error("error while fetching all the records error is ",err)
		return nil,0, GetErrorResponseByCode(InternalServerError)
	}
	return books,totalCount, nil
}

func (s *bookService) CreateOrUpdate(ctx context.Context, id int, req CreateOrUpdateBookRequest) (int64, *ErrorResponse) {
	b, err := s.Get(ctx, id)
	if err != nil {
		if err == GetErrorResponseByCode(BookNotFound) {
			logrus.Info("no record exist for given id  ",id," creating the record")
			id, createErr := s.Create(ctx, req)
			if createErr != nil {
				logrus.Error("error while fetching creating the record. error is ",err)
				return 0, GetErrorResponseByCode(InternalServerError)
			}
			return id, nil
		}
		return 0, err
	}
	if req.Title != "" {
		b.Title = req.Title
	}
	if req.Author != "" {
		b.Author = req.Author
	}
	if req.Description != "" {
		b.Description = req.Description
	}
	if err := s.repository.Update(ctx, b); err != nil {
		logrus.Error("error while updating the record. error is ",err)
		return 0, GetErrorResponseByCode(InternalServerError)
	}
	return 0, nil
}

func (s *bookService) Delete(ctx context.Context, id int) *ErrorResponse {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return GetErrorResponseByCode(InternalServerError)
	}
	return nil
}
