package dto

type GameListCreateRequest struct {
	Name     string `json:"name" validate:"required;min=3;max=100"`
	IsPublic bool   `json:"isPublic"`
}

type GameListUpdateRequest struct {
	Name     *string `json:"name" validate:"min=3;max=100"`
	IsPublic *bool   `json:"isPublic"`
}

type GameListResponse struct {
	Name string `json:"name"`
}
