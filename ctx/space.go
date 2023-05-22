package ctx

import (
	"context"
	"net/http"

	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/httperror"
)

var spaceKey = contextKey("space")

func GetSpace(ctx context.Context) *models.Space {
	return ctx.Value(spaceKey).(*models.Space)
}

func SetSpace(spacesRepo repositories.SpacesRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r.Context())

			spaceID := c.PathParam("spaceID")

			space, err := spacesRepo.Get(r.Context(), spaceID)
			if err != nil {
				c.HandleError(err)
				return
			}

			ctx := context.WithValue(r.Context(), spaceKey, space)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CanViewSpace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := Get(r.Context())
		space := GetSpace(r.Context())

		if !c.IsSpaceAdmin(c.User, space) {
			c.HandleError(httperror.Forbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
