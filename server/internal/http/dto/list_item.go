package dto

import "github.com/google/uuid"

type ListItemAddRequest struct {
	Status string    `json:"status"`
	Rating float32   `json:"rating"`
	GameID uuid.UUID `query:"game_id"`
	ListID uuid.UUID `query:"list_id"`
}

type ListItemUpdateRequest struct {
	Status *string   `json:"status,omitempty"`
	Rating *float32  `json:"rating,omitempty"`
	GameID uuid.UUID `json:"game_id"`
	ListID uuid.UUID `query:"list_id"`
}

type ListItemDeleteRequest struct {
	GameID uuid.UUID `query:"game_id"`
	ListID uuid.UUID `query:"list_id"`
}

type ListItemResponse struct {
	Status string    `json:"status"`
	Rating float32   `json:"rating"`
	GameID uuid.UUID `query:"game_id"`
}
