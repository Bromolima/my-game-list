package dto

type GameCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Genre       string `json:"genre" validate:"required,min=3,max=100"`
	Developer   string `json:"developer" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=10"`
	ImageURL    string `json:"image_url,omitempty" validate:"url"`
}

type GameUpdateRequest struct {
	Name        *string `json:"name,omitempty" validate:"min=3,max=100"`
	Genre       *string `json:"genre,omitempty" validate:"min=3,max=100"`
	Developer   *string `json:"developer,omitempty" validate:"min=3,max=100"`
	Description *string `json:"description" validate:"min=10"`
	ImageURL    *string `json:"image_url,omitempty" validate:"url"`
}

type GamesSearchRequest struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Name  string `query:"name"`
}

type GameResponse struct {
	Name        string  `json:"name"`
	Genre       string  `json:"genre"`
	Developer   string  `json:"developer"`
	Rating      float32 `json:"rating"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
}
