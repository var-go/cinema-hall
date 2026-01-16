package dto

type CreateUserRequest struct {
	Email    string `json:"email" binding:"request,email"`
	Password string `json:"password" binding:"request,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Email *string `json:"email"`
	Name  *string `json:"name"`
	Role  *string `json:"role"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}
