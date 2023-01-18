package models

import (
	"time"

	"github.com/ugent-library/deliver/validate"
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
	return validate.Validate(
		validate.NotEmpty("name", s.Name),
		validate.LengthIn("name", s.Name, 1, 50),
		validate.Alphanumeric("name", s.Name),
	)
}
