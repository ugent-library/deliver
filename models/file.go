package models

import (
	"time"

	"github.com/ugent-library/okay"
)

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

func (f *File) Validate() error {
	return okay.Validate(
		okay.NotEmpty("name", f.Name),
		okay.Min("size", f.Size, 1),
	)
}
