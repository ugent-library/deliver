package ctx

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/repository"
	"github.com/ugent-library/httperror"
)

var spaceKey = contextKey("space")

func GetSpace(r *http.Request) *models.Space {
	return r.Context().Value(spaceKey).(*models.Space)
}

func SetSpace(spacesRepo repository.SpacesRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			spaceName := chi.URLParam(r, "spaceName")

			space, err := spacesRepo.GetByName(r.Context(), spaceName)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), spaceKey, space)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CanViewSpace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := Get(r)
		space := GetSpace(r)

		if !c.Permissions.IsSpaceAdmin(c.User, space) {
			c.HandleError(w, r, httperror.Forbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
