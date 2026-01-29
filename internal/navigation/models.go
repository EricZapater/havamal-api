package navigation

import "github.com/google/uuid"

type Type string
type LinkSource string

const (
	TypeInternal Type = "internal"
	TypeExternal Type = "external"
	
	LinkSourceCustom   LinkSource = "custom"
	LinkSourceCategory LinkSource = "category"
	LinkSourcePost     LinkSource = "post"
)

type Request struct {
	Label      string     `json:"label"`
	Slug       string     `json:"slug"`
	Type       Type       `json:"type"`
	Order      int        `json:"order"`
	ParentId   string     `json:"parent_id"`
	LinkSource LinkSource `json:"link_source"`
	CategoryId string     `json:"category_id"`
	PostId     string     `json:"post_id"`
}

type Navigation struct {
	ID         uuid.UUID   `json:"id"`
	Label      string      `json:"label"`
	Slug       string      `json:"slug"`
	Type       Type        `json:"type"`
	Order      int         `json:"order"`
	ParentId   *uuid.UUID  `json:"parent_id"`
	LinkSource LinkSource  `json:"link_source"`
	CategoryId *uuid.UUID  `json:"category_id"`
	PostId     *uuid.UUID  `json:"post_id"`
}