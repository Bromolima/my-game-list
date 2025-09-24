package dto

type UserRegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,containsany=@#!&$*"`
	Username  string `json:"username" validate:"required,min=6,max=20"`
	AvatarURL string `json:"avatar_url" validate:"url"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,containsany=@#!&$*"`
}

type UserUpdateRequest struct {
	Email     *string `json:"email,omitempty" validate:"email"`
	Password  *string `json:"password,omitempty" validate:"min=6,containsany=@#!&$*"`
	Username  *string `json:"username,omitempty" validate:"min=6,max=20"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"url"`
}

type UserSearchRequest struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Name  string `query:"name"`
}

type UserResponse struct {
	Username string `json:"username"`
}
