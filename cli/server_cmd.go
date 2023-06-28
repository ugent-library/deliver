package cli

import (
	"context"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/nics/ich"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/handlers"
	"github.com/ugent-library/deliver/htmx"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstore"
	"github.com/ugent-library/deliver/repositories"
	"github.com/ugent-library/httpx/render"
	mw "github.com/ugent-library/middleware"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// setup services
		repo, err := repositories.New(config.Repo.Conn)
		if err != nil {
			return err
		}

		storage, err := objectstore.New(config.Storage.Backend, config.Storage.Conn)
		if err != nil {
			return err
		}

		// setup permissions
		permissions := &models.Permissions{
			Admins: config.Admins,
		}

		// setup auth
		oidcAuth, err := oidc.NewAuth(context.Background(), oidc.Config{
			URL:          config.OIDC.URL,
			ClientID:     config.OIDC.ID,
			ClientSecret: config.OIDC.Secret,
			RedirectURL:  config.OIDC.RedirectURL,
			CookieName:   "deliver.state",
			CookieSecret: []byte(config.Cookie.Secret),
		})
		if err != nil {
			return err
		}

		// setup health checker
		// TODO add checkers
		healthChecker := health.NewChecker()

		// setup assets
		assets, err := mix.New(mix.Config{
			ManifestFile: "static/mix-manifest.json",
			PublicPath:   "/static/",
		})
		if err != nil {
			return err
		}

		// setup router
		router := ich.New()
		router.Use(middleware.RequestID)
		if config.Env != "local" {
			router.Use(middleware.RealIP)
		}
		router.Use(mw.MethodOverride( // TODO eliminate need for method override
			mw.MethodFromHeader(mw.MethodHeader),
			mw.MethodFromForm(mw.MethodParam),
		))
		router.Use(zaphttp.SetLogger(logger.Desugar(), zapchi.RequestID))
		router.Use(middleware.RequestLogger(zapchi.LogFormatter()))
		router.Use(middleware.Recoverer)
		router.Use(middleware.StripSlashes)

		// htmx message hub
		hub := htmx.NewHub(htmx.Config{
			// TODO htmx secret config
			Secret: []byte(config.Cookie.Secret),
		})

		// routes
		router.Get("/health", health.NewHandler(healthChecker))
		router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, http.StatusOK, &struct {
				Branch string `json:"branch,omitempty"`
				Commit string `json:"commit,omitempty"`
			}{
				Branch: config.Source.Branch,
				Commit: config.Source.Commit,
			})
		})
		router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
			if err := hub.HandleWebSocket(w, r, r.URL.Query().Get("channels")); err != nil {
				logger.Error(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		})
		router.Group(func(r *ich.Mux) {
			r.Use(
				func(next http.Handler) http.Handler {
					return http.MaxBytesHandler(next, config.MaxFileSize)
				},
				csrf.Protect(
					[]byte(config.Cookie.Secret),
					csrf.CookieName("deliver.csrf"),
					csrf.Path("/"),
					csrf.HttpOnly(true),
					csrf.SameSite(csrf.SameSiteStrictMode),
					csrf.FieldName("_csrf_token"),
				),
				// request context wrapper
				ctx.Set(ctx.Config{
					Repo:        repo,
					Storage:     storage,
					MaxFileSize: config.MaxFileSize,
					Auth:        oidcAuth,
					Router:      router,
					ErrorHandlers: map[int]http.HandlerFunc{
						http.StatusNotFound:     handlers.NotFound,
						http.StatusUnauthorized: handlers.Unauthorized,
						http.StatusForbidden:    handlers.Forbidden,
					},
					Permissions: permissions,
					Assets:      assets,
					Hub:         hub,
					Env:         config.Env,
				}),
			)

			// viewable by everyone
			r.NotFound(handlers.NotFound)
			r.Get("/", handlers.Home).Name("home")
			r.Get("/auth/callback", handlers.AuthCallback)
			r.Get("/login", handlers.Login).Name("login")
			r.Get("/logout", handlers.Logout).Name("logout")
			r.With(ctx.SetFolder(*repo.Folders)).Get("/share/{folderID}:{folderSlug}", handlers.ShareFolder).Name("shareFolder")
			r.With(ctx.SetFile(*repo.Files)).Get("/files/{fileID}", handlers.DownloadFile).Name("downloadFile")
			// viewable by space owners and admins
			r.Group(func(r *ich.Mux) {
				r.Use(ctx.RequireUser)
				r.Get("/spaces", handlers.ListSpaces).Name("spaces")
				r.With(ctx.RequireAdmin).Get("/new-space", handlers.NewSpace).Name("newSpace")
				r.With(ctx.RequireAdmin).Post("/spaces", handlers.CreateSpace).Name("createSpace")
				r.Route("/spaces/{spaceName}", func(r *ich.Mux) {
					r.Use(ctx.SetSpace(*repo.Spaces))
					r.Use(ctx.CanViewSpace)
					r.Get("/", handlers.ShowSpace).Name("space")
					r.Post("/folders", handlers.CreateFolder).Name("createFolder")
					r.With(ctx.RequireAdmin).Get("/edit", handlers.EditSpace).Name("editSpace")
					r.With(ctx.RequireAdmin).Put("/", handlers.UpdateSpace).Name("updateSpace")
				})
				r.Route("/folders/{folderID}", func(r *ich.Mux) {
					r.Use(ctx.SetFolder(*repo.Folders))
					r.Use(ctx.CanEditFolder)
					r.Get("/", handlers.ShowFolder).Name("folder")
					r.Get("/edit", handlers.EditFolder).Name("editFolder")
					r.Put("/", handlers.UpdateFolder).Name("updateFolder")
					r.Post("/files", handlers.UploadFile).Name("uploadFile")
					r.Delete("/", handlers.DeleteFolder).Name("deleteFolder")
				})
				r.With(ctx.SetFile(*repo.Files), ctx.CanEditFile).Delete("/files/{fileID}", handlers.DeleteFile).Name("deleteFile")
			})
		})

		// start server
		server := graceful.WithDefaults(&http.Server{
			Addr:         config.Addr(),
			Handler:      router,
			ReadTimeout:  10 * time.Minute,
			WriteTimeout: 10 * time.Minute,
		})
		logger.Infof("starting server at %s", config.Addr())
		if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
			return err
		}
		logger.Info("gracefully stopped server")

		return nil
	},
}
