package integrationtest

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestCreateBook_ShouldBeSuccessful(t *testing.T) {
	req := Request{
		URL:                    "/books",
		MethodType:             "POST",
		RequestBodyFilePath:    "./request/create_book_request.json",
		ExpectedHttpStatusCode: http.StatusCreated,
		ExpectedHeaders: map[string]string{
			"Location": "/books/1",
		},
	}
	Exec(t, req)
}

func TestGetBookById_ShouldBeSuccessful(t *testing.T) {
	id, err := insertTestBook(t.Context(), "Test", "Test Author", "Test Desc")
	if err != nil {
		logrus.Fatalf("error while inserting data in db %s", err)
	}
	req := Request{
		URL:                          "/books/" + strconv.Itoa(int(id)),
		MethodType:                   "GET",
		ExpectedHttpStatusCode:       http.StatusOK,
		ExpectedResponseBodyFilePath: "response/get_book_response.json",
	}
	Exec(t, req)
}

func TestGetAll_ShouldBeSuccessful(t *testing.T) {
	_, err := insertTestBook(t.Context(), "Test", "Test Author", "Test Desc")
	if err != nil {
		logrus.Fatalf("error while inserting data in db %s", err)
	}
	_, err = insertTestBook(t.Context(), "Test2", "Test Author2", "Test Desc2")
	if err != nil {
		logrus.Fatalf("error while inserting data in db %s", err)
	}
	req := Request{
		URL:                          "/books?page=1&limit=10",
		MethodType:                   "GET",
		ExpectedHttpStatusCode:       http.StatusOK,
		ExpectedResponseBodyFilePath: "response/get_paginated_book_response.json",
	}
	Exec(t, req)
}

func TestPut_ShouldUpdateTheExistingBook(t *testing.T) {
	id, err := insertTestBook(t.Context(), "Test", "Test Author", "Test Desc")
	if err != nil {
		logrus.Fatalf("error while inserting data in db %s", err)
	}
	req := Request{
		URL:                    "/books/" + strconv.Itoa(int(id)),
		MethodType:             "PUT",
		RequestBodyFilePath:    "./request/update_book_request.json",
		ExpectedHttpStatusCode: http.StatusNoContent,
	}
	Exec(t, req)
}

func TestPut_ShouldCreateNewBookWhenBookWithCurrentIdDoesnotExist(t *testing.T) {
	req := Request{
		URL:                    "/books/" + strconv.Itoa(int(100)),
		MethodType:             "PUT",
		RequestBodyFilePath:    "./request/update_book_request.json",
		ExpectedHttpStatusCode: http.StatusNoContent,
		ExpectedHeaders: map[string]string{
			"Location": "/books/1",
		},
	}
	Exec(t, req)
}

func TestDelete_ShouldBeSuccessful(t *testing.T) {
	id, err := insertTestBook(t.Context(), "Test", "Test Author", "Test Desc")
	if err != nil {
		logrus.Fatalf("error while inserting data in db %s", err)
	}
	req := Request{
		URL:                    "/books/" + strconv.Itoa(int(id)),
		MethodType:             "DELETE",
		ExpectedHttpStatusCode: http.StatusNoContent,
	}
	Exec(t, req)
}
