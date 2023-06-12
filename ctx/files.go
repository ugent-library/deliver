package ctx

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/httperror"
)

var fileKey = contextKey("file")

func GetFile(r *http.Request) *models.File {
	return r.Context().Value(fileKey).(*models.File)
}

func SetFile(filesRepo repositories.FilesRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)

			fileID := chi.URLParam(r, "fileID")

			file, err := filesRepo.Get(r.Context(), fileID)
			if err != nil {
				c.HandleError(w, r, err)
				return
			}

			ctx := context.WithValue(r.Context(), fileKey, file)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CanEditFile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := Get(r)
		file := GetFile(r)

		if !c.IsSpaceAdmin(c.User, file.Folder.Space) {
			c.HandleError(w, r, httperror.Forbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
