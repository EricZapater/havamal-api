package versions

import (
	"time"

	"github.com/google/uuid"
)

type Request struct {
	Version       string    `json:"version"`
	PostId        string    `json:"post_id"`
	VersionNumber int       `json:"version_number"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
}

type Version struct {
	ID            uuid.UUID `json:"id"`
	Version       string    `json:"version"`
	PostId        uuid.UUID `json:"post_id"`
	VersionNumber int       `json:"version_number"`
	Content       string    `json:"content"`
	CreatedAt     time.Time `json:"created_at"`
}
