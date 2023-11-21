package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/ugent-library/okay"
)

type User struct {
	ID            string    `json:"id,omitempty"`
	Username      string    `json:"username,omitempty"`
	Name          string    `json:"name,omitempty"`
	Email         string    `json:"email,omitempty"`
	RememberToken string    `json:"remember_token,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

func (u *User) Validate() error {
	return okay.Validate(
		okay.NotEmpty("username", u.Username),
		okay.LengthIn("username", u.Username, 1, 50),
		okay.NotEmpty("name", u.Name),
		okay.NotEmpty("email", u.Email),
	)
}

func NewRememberToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
