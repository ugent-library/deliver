package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/nics/ich"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	"github.com/ugent-library/crypt"
	"github.com/ugent-library/deliver/catbird"
	"github.com/ugent-library/httpx"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"github.com/ugent-library/zaphttp/zapchi"
	"github.com/unrolled/secure"
	"github.com/unrolled/secure/cspbuilder"

	"github.com/ugent-library/deliver/ctx"
	"github.com/ugent-library/deliver/handlers"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/objectstores"
	"github.com/ugent-library/deliver/repositories"
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

		storage, err := objectstores.New(config.Storage.Backend, config.Storage.Conn)
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
			CookiePrefix: "deliver.",
		})
		if err != nil {
			return err
		}

		// setup assets
		assets, err := loadAssets("static/manifest.json")
		if err != nil {
			return err
		}

		// setup htmx message hub
		hub, err := catbird.New(catbird.Config{})
		if err != nil {
			return err
		}
		defer hub.Stop()

		// setup timezone
		timezone, err := time.LoadLocation(config.Timezone)
		if err != nil {
			return err
		}

		// setup handlers
		authHandler := handlers.NewAuthHandler(
			oidcAuth,
			config.OIDC.UsernameClaim,
			config.OIDC.NameClaim,
			config.OIDC.EmailClaim,
		)

		// setup router
		router := ich.New()
		router.Use(middleware.RequestID)
		if config.Env != "local" {
			router.Use(middleware.RealIP)
		}
		router.Use(httpx.MethodOverride)
		router.Use(zaphttp.SetLogger(logger.Desugar(), zapchi.RequestID))
		router.Use(middleware.RequestLogger(zapchi.LogFormatter()))
		router.Use(middleware.Recoverer)
		router.Use(middleware.StripSlashes)
		router.Use(secure.New(secure.Options{
			IsDevelopment: config.Env == "local",
			ContentSecurityPolicy: (&cspbuilder.Builder{
				Directives: map[string][]string{
					cspbuilder.DefaultSrc: {"'self'"},
					cspbuilder.ScriptSrc:  {"'self'", "$NONCE", "blob:"},
					cspbuilder.StyleSrc:   {"'self'", "'unsafe-inline'"}, // htmx injects style tags
					cspbuilder.ImgSrc:     {"'self'", "data:"},           // bootstrap uses data: images
				},
			}).MustBuild(),
		}).Handler)

		// mount health and info
		router.Get("/status", health.NewHandler(health.NewChecker())) // TODO add checkers
		router.Get("/info", func(w http.ResponseWriter, r *http.Request) {
			httpx.RenderJSON(w, http.StatusOK, version)
		})

		// mount assets
		router.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

		// mount ui routes
		router.Group(func(r *ich.Mux) {
			// TODO use new RequestSize middleware
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
					csrf.Secure(config.Env != "local"),
				),
				// request context wrapper
				ctx.Set(ctx.Config{
					Crypt:       crypt.New([]byte(config.Cookie.Secret)),
					Repo:        repo,
					Storage:     storage,
					MaxFileSize: config.MaxFileSize,
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
					Timezone:    timezone,
					CSRFName:    "_csrf_token",
				}),
			)

			// viewable by everyone
			r.NotFound(handlers.NotFound)
			r.Get("/", handlers.Home).Name("home")
			r.Get("/auth/callback", authHandler.AuthCallback)
			r.Get("/login", authHandler.Login).Name("login")
			r.Get("/logout", authHandler.Logout).Name("logout")
			r.With(ctx.SetFolder(*repo.Folders)).Get("/share/{folderID}:{folderSlug}", handlers.ShareFolder).Name("shareFolder")
			r.With(ctx.SetFile(*repo.Files)).Get("/files/{fileID}", handlers.DownloadFile).Name("downloadFile")
			r.With(ctx.SetFolder(*repo.Folders)).Get("/folders/{folderID}.zip", handlers.DownloadFolder).Name("downloadFolder")
			// viewable by space owners and admins
			r.Group(func(r *ich.Mux) {
				r.Use(ctx.RequireUser)

				// mount htmx message hub
				r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
					c := ctx.Get(r)
					var topics []string
					if err := c.DecryptValue(r.URL.Query().Get("token"), &topics); err != nil {
						c.HandleError(w, r, err)
						return
					}
					if err := hub.HandleWebsocket(w, r, c.User.ID, topics); err != nil {
						c.HandleError(w, r, err)
						return
					}
				}).Name("ws")

				r.Get("/spaces", handlers.ListSpaces).Name("spaces")
				r.With(ctx.RequireAdmin).Get("/new-space", handlers.NewSpace).Name("newSpace")
				r.With(ctx.RequireAdmin).Post("/spaces", handlers.CreateSpace).Name("createSpace")
				r.Route("/spaces/{spaceName}", func(r *ich.Mux) {
					r.Use(ctx.SetSpace(*repo.Spaces))
					r.Use(ctx.CanViewSpace)
					r.Get("/", handlers.ShowSpace).Name("space")
					r.Get("/folders", handlers.GetFolders).Name("getFolders")
					r.Post("/", handlers.CreateFolder).Name("createFolder")
					r.With(ctx.RequireAdmin).Get("/edit", handlers.EditSpace).Name("editSpace")
					r.With(ctx.RequireAdmin).Put("/", handlers.UpdateSpace).Name("updateSpace")
				})
				r.Route("/folders/{folderID}", func(r *ich.Mux) {
					r.Use(ctx.SetFolder(*repo.Folders))
					r.Use(ctx.CanEditFolder)
					r.Get("/", handlers.ShowFolder).Name("folder")
					r.Get("/edit", handlers.EditFolder).Name("editFolder")
					r.Put("/", handlers.UpdateFolder).Name("updateFolder")
					r.Put("/postpone", handlers.PostponeFolderExpiration).Name("postponeExpiration")
					r.Post("/files", handlers.UploadFile).Name("uploadFile")
					r.Delete("/", handlers.DeleteFolder).Name("deleteFolder")
				})
				r.With(ctx.SetFile(*repo.Files), ctx.CanEditFile).
					Delete("/files/{fileID}", handlers.DeleteFile).
					Name("deleteFile")
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

func loadAssets(manifestFile string) (map[string]string, error) {
	data, err := os.ReadFile(manifestFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't read manifest '%s': %w", manifestFile, err)
	}

	assets := make(map[string]string)
	if err = json.Unmarshal(data, &assets); err != nil {
		return nil, fmt.Errorf("couldn't parse manifest '%s': %w", manifestFile, err)
	}

	return assets, nil
}
