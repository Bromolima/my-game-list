package dto

type UserCreatePayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,containsany=@#!&$*"`
	Username string `json:"username" validate:"required,min=6,max=20"`
}

type UserLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,containsany=@#!&$*"`
}

type UserResponse struct {
	Username string `json:"username"`
}
