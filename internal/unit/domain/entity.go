package domain

import (
	"time"
)

type Unit struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Symbol    string    `db:"symbol" json:"symbol"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateUnitRequest struct {
	Name   string `json:"name" validate:"required,min=1"`
	Symbol string `json:"symbol" validate:"required,min=1"`
}

type UpdateUnitRequest struct {
	Name   string `json:"name" validate:"required,min=1"`
	Symbol string `json:"symbol" validate:"required,min=1"`
}

type ListUnitsRequest struct {
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	Search string `json:"search,omitempty"`
}

type UnitResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	CreatedAt time.Time `json:"created_at"`
}

type UnitListResponse struct {
	Units []UnitResponse `json:"units"`
	Meta  ListMeta       `json:"meta"`
}

type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
