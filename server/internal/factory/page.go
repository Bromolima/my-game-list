package factory

import (
	"github.com/Bromolima/my-game-list/internal/entities"
	"github.com/Bromolima/my-game-list/internal/http/dto"
)

const (
	DefaultLimit = 10
	MaxLimit     = 50
	DefaultPage  = 1
)

func NewPage[T any](page, limit int) *entities.Page[T] {
	if page <= 0 {
		page = DefaultPage
	}

	if limit <= 0 {
		limit = DefaultLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return &entities.Page[T]{
		Data:   make([]T, 0),
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func NewReponseFromPage[model any, response any](p *entities.Page[model], toResponse func(m *model) *response) *dto.PageResponse[response] {
	pageResponse := &dto.PageResponse[response]{
		Page:       p.Page,
		Limit:      p.Limit,
		TotalPages: p.TotalPages,
		Data:       make([]response, 0),
	}

	for _, data := range p.Data {
		resp := toResponse(&data)
		pageResponse.Data = append(pageResponse.Data, *resp)
	}

	return pageResponse
}
