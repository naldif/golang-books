package books

import (
	"time"

	"golang-buku/internal/modules/categories"
)

type CreateBookInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ReleaseYear int    `json:"release_year" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	TotalPage   int    `json:"total_page" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
}

type UpdateBookInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ReleaseYear int    `json:"release_year" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	TotalPage   int    `json:"total_page" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
}

type BookResponse struct {
	ID          uint                 `json:"id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	ImageURL    string               `json:"image_url"`
	ReleaseYear int                  `json:"release_year"`
	Price       int                  `json:"price"`
	TotalPage   int                  `json:"total_page"`
	Thickness   string               `json:"thickness"`
	CategoryID  uint                 `json:"category_id"`
	Category    *categories.Category `json:"category,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	CreatedBy   string               `json:"created_by"`
	ModifiedAt  time.Time            `json:"modified_at"`
	ModifiedBy  string               `json:"modified_by"`
}
