package book

import "net/http"

type ErrorResponse struct {
	HttpStatusCode int       `json:"-"`
	ErrorCode      ErrorCode `json:"errorCode" example:"BAD_REQUEST"`
	ErrorMessage   string    `json:"errorMessage" example:"limit must be >=1"`
}

func (e ErrorResponse) Error() string {
	return e.ErrorMessage
}

var errorResponseMap = map[ErrorCode]*ErrorResponse{
	BookNotFound: {
		HttpStatusCode: http.StatusNotFound,
		ErrorCode:      BookNotFound,
		ErrorMessage:   "book not found",
	},
	InternalServerError: {
		HttpStatusCode: http.StatusInternalServerError,
		ErrorCode:      InternalServerError,
		ErrorMessage:   "internal server error",
	},
	BadRequest: {
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      BadRequest,
		ErrorMessage:   "request is invalid.",
	},
}

func GetErrorResponseByCode(errCode ErrorCode) *ErrorResponse {
	return errorResponseMap[errCode]
}

func GetErrorResponse(errCode ErrorCode,errorMessage string,statusCode int) *ErrorResponse {
	return &ErrorResponse{
		HttpStatusCode: statusCode,
		ErrorCode:      errCode,
		ErrorMessage:   errorMessage,
	}
}

type ErrorCode string

const (
	BookNotFound        ErrorCode = "BOOK_NOT_FOUND"
	InternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	BadRequest          ErrorCode = "BAD_REQUEST"
)
