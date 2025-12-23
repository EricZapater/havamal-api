package posts

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Draft     Status = "draft"
	Published Status = "published"
	Archived  Status = "archived"
)

type Request struct {
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"`
	Status      Status    `json:"status"`
	PublishedAt time.Time `json:"published_at"`	
	UpdatedAt   time.Time `json:"updated_at"`
	AuthorId 	string `json:"author_id"`
	Author      string `json:"author"`
	CategoryId  string `json:"categoryId"`
}

type Post struct {
	ID 			uuid.UUID `json:"id"`
	Title 		string `json:"title"`
	Slug 		string `json:"slug"`
	Summary 	string `json:"summary"`
	Content 	string `json:"content"`
	Status 		Status `json:"status"`
	PublishedAt 	time.Time `json:"published_at"`	
	UpdatedAt 	time.Time `json:"updated_at"`
	AuthorId 	uuid.UUID `json:"author_id"`
}

type Response struct{
	ID 	uuid.UUID `json:"id"`
	Title 	string `json:"title"`
	Slug 	string `json:"slug"`
	Summary 	string `json:"summary"`
	Content 	string `json:"content"`
	Status 	Status `json:"status"`
	PublishedAt 	time.Time `json:"published_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	AuthorId 	uuid.UUID `json:"author_id"`
	CategoryId uuid.UUID `json:"category_id"`
	CategoryName string `json:"category_name"`
	CategoryDescription string `json:"category_description"`
	CategorySlug string `json:"category_slug"`
	AuthorName string `json:"author_name"`

}

type PostCategories struct {
	PostId 	uuid.UUID `json:"post_id"`
	CategoryId 	uuid.UUID `json:"category_id"`
}

type PostVersion struct {	
	VersionId 	uuid.UUID `json:"version"`
	PostId 	uuid.UUID `json:"post_id"`
}