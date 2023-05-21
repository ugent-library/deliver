package ctx

import (
	"context"
	"net/http"

	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/httperror"
)

var fileKey = contextKey("file")

func GetFile(ctx context.Context) *models.File {
	return ctx.Value(fileKey).(*models.File)
}

func SetFile(filesRepo repositories.FilesRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r.Context())

			fileID := c.PathParam("fileID")

			file, err := filesRepo.Get(r.Context(), fileID)
			if err != nil {
				c.HandleError(err)
				return
			}

			ctx := context.WithValue(r.Context(), fileKey, file)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CanEditFile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := Get(r.Context())
		file := GetFile(r.Context())

		if !c.IsSpaceAdmin(c.User, file.Folder.Space) {
			c.HandleError(httperror.Forbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
