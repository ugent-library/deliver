package controllers

import "github.com/ugent-library/dilliver/httperror"

func RequireUser(c Ctx) error {
	if c.User() == nil {
		return httperror.Unauthorized
	}
	return nil
}
