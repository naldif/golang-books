package categories

import "time"

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

type CategoryResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy string    `json:"modified_by"`
}
