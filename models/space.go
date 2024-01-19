package models

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/ugent-library/okay"
)

type Space struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Admins    []string  `json:"admins,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// relations (can be empty)
	Folders []*Folder `json:"folders,omitempty"`
}

func (s *Space) Validate() error {
	return okay.Validate(
		okay.NotEmpty("name", s.Name),
		okay.LengthBetween("name", s.Name, 1, 50),
		okay.Alphanumeric("name", s.Name),
	)
}

func (s *Space) Fake(faker *gofakeit.Faker) (any, error) {
	return Space{
		Name:   gofakeit.RandomString([]string{"BIBXYZ", "ABCLIB", "DEFCOL", "UNIZXY", "department", "BIBLIB", "FACLIB"}),
		Admins: []string{"deliver"},
	}, nil
}
