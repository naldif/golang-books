package users

import "time"

type CreateUserInput struct {
	Username string `json:"username" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  string    `json:"created_by"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy string    `json:"modified_by"`
}
