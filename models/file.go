package models

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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

func (f *File) Fake(faker *gofakeit.Faker) (any, error) {
	return File{
		Name:      fmt.Sprintf("%d.jpg", faker.Number(123456, 912345)),
		Downloads: int64(faker.Number(0, 10)),
	}, nil
}
