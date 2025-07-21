package book

type CreateOrUpdateBookRequest struct {
    Title       string `json:"title" validate:"required,min=1,max=200"`
    Author      string `json:"author" validate:"required,min=1,max=100"`
    Description string `json:"description" validate:"omitempty,max=500"`
}

// Response DTO
type BookResponse struct {
	ID          int    `json:"id" example:"1"`
	Title       string `json:"title" example:"Harry Potter"`
	Author      string `json:"author" example:"JK Rolling"`
	Description string `json:"description" example:"harry potter and his friends"`
}

type PaginatedBookListResponse struct {
  Page       int            `json:"page" example:"1"`
  Limit      int            `json:"limit" example:"10"`
  Total      int            `json:"total" example:"42"`
  TotalPages int            `json:"totalPages" example:"5"`
  Data       []BookResponse `json:"data"`
}
