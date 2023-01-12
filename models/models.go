package models

import (
	"time"

	"github.com/ugent-library/deliver/validate"
)

type Space struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// relations (can be empty)
	Folders []*Folder `json:"folders,omitempty"`
}

type Folder struct {
	ID        string    `json:"id,omitempty"`
	SpaceID   string    `json:"space_id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	// relations (can be empty)
	Size      int64   `json:"size"`
	FileCount int     `json:"file_count"`
	Space     *Space  `json:"space,omitempty"`
	Files     []*File `json:"files,omitempty"`
}

type File struct {
	ID          string    `json:"id,omitempty"`
	FolderID    string    `json:"folder_id,omitempty"`
	MD5         string    `json:"md5,omitempty"`
	Name        string    `json:"name,omitempty"`
	Size        int64     `json:"size,omitempty"`
	ContentType string    `json:"content_type,omitempty"`
	Downloads   int64     `json:"downloads,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	// relations (can be empty)
	Folder *Folder `json:"folder,omitempty"`
}

func (s *Space) Validate() error {
	return validate.Validate(
		validate.NotEmpty("name", s.Name),
		validate.LengthIn("name", s.Name, 1, 256),
	)
}

func (f *Folder) Validate() error {
	return validate.Validate(
		validate.NotEmpty("name", f.Name),
		validate.LengthIn("name", f.Name, 1, 256),
	)
}

func (f *File) Validate() error {
	return validate.Validate(
		validate.NotEmpty("name", f.Name),
		validate.Min("size", f.Size, 1),
	)
}
