package validators

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Email string `json:"email" binding:"required,email,max=100"`
	Phone string `json:"phone" binding:"max=20"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Email string `json:"email" binding:"required,email,max=100"`
	Phone string `json:"phone" binding:"max=20"`
}
