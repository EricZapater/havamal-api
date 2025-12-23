package navigation

import "github.com/google/uuid"

type Type string

const (
	TypeInternal Type = "internal"
	TypeExternal Type = "external"
)

type Request struct {
	Label    string `json:"label"`
	Slug     string `json:"slug"`
	Type     Type   `json:"type"`
	Order    int    `json:"order"`
	ParentId string `json:"parent_id"`
}

type Navigation struct {
	ID       uuid.UUID `json:"id"`
	Label    string    `json:"label"`
	Slug     string    `json:"slug"`
	Type     Type      `json:"type"`
	Order    int       `json:"order"`
	ParentId *uuid.UUID `json:"parent_id"`
}