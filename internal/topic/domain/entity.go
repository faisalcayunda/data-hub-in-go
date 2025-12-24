package domain

import (
	"time"
)

type Topic struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateTopicRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type UpdateTopicRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type ListTopicsRequest struct {
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	Search string `json:"search,omitempty"`
}

type TopicResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type TopicListResponse struct {
	Topics []TopicResponse `json:"topics"`
	Meta   ListMeta        `json:"meta"`
}

type ListMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
