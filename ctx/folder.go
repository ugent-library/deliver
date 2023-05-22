package ctx

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/httperror"
)

var folderKey = contextKey("folder")

func GetFolder(ctx context.Context) *models.Folder {
	return ctx.Value(folderKey).(*models.Folder)
}

func SetFolder(foldersRepo repositories.FoldersRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r.Context())

			folderID := chi.URLParam(r, "folderID")

			folder, err := foldersRepo.Get(r.Context(), folderID)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), folderKey, folder)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CanEditFolder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := Get(r.Context())
		folder := GetFolder(r.Context())

		if !c.IsSpaceAdmin(c.User, folder.Space) {
			c.HandleError(w, r, httperror.Forbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
