//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=views
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/oklog/ulid/v2"
	"github.com/ory/graceful"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	c "github.com/ugent-library/deliver/controllers"
	"github.com/ugent-library/deliver/controllers/ctx"
	"github.com/ugent-library/deliver/crumb"
	"github.com/ugent-library/deliver/models"
	"github.com/ugent-library/deliver/turbo"
	"github.com/ugent-library/deliver/turborouter"
	"github.com/ugent-library/friendly"
	"github.com/ugent-library/middleware"
	"github.com/ugent-library/mix"
	"github.com/ugent-library/oidc"
	"github.com/ugent-library/zaphttp"
	"github.com/unrolled/render"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(appCmd)
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Start the web app server",
	Run: func(cmd *cobra.Command, args []string) {
		// setup services
		repoService, err := models.NewRepositoryService(models.RepositoryConfig{
			DB: config.DB,
		})
		if err != nil {
			logger.Fatal(err)
		}

		fileService, err := models.NewFileService(models.FileConfig{
			S3URL:    config.S3.URL,
			S3ID:     config.S3.ID,
			S3Secret: config.S3.Secret,
			S3Bucket: config.S3.Bucket,
			S3Region: config.S3.Region,
		})
		if err != nil {
			logger.Fatal(err)
		}

		// setup permissions
		permissions := &models.Permissions{
			Admins: config.Admins,
		}

		// setup auth
		oidcAuth, err := oidc.NewAuth(context.TODO(), oidc.Config{
			URL:          config.OIDC.URL,
			ClientID:     config.OIDC.ID,
			ClientSecret: config.OIDC.Secret,
			RedirectURL:  config.OIDC.RedirectURL,
			CookieName:   "deliver.state",
			CookieSecret: []byte(config.CookieSecret),
		})
		if err != nil {
			logger.Fatal(err)
		}

		// setup assets
		assets, err := mix.New(mix.Config{
			ManifestFile: "static/mix-manifest.json",
			PublicPath:   "/static/",
		})
		if err != nil {
			logger.Fatal(err)
		}

		// setup renderer
		renderer := render.New(render.Options{
			Directory:          "templates",
			Extensions:         []string{".gohtml"},
			RequirePartials:    true,
			HTMLTemplateOption: "missingkey=error",
			Funcs: []template.FuncMap{{
				"assetPath":     assets.AssetPath,
				"friendlyBytes": friendly.Bytes,
				"join":          strings.Join,
			}},
		})

		// setup router
		r := mux.NewRouter()
		r.StrictSlash(true)
		r.UseEncodedPath()

		// controllers
		errs := c.NewErrors()
		auth := c.NewAuth(repoService, oidcAuth)
		pages := c.NewPages()
		spaces := c.NewSpaces(repoService)
		folders := c.NewFolders(repoService, fileService, viper.GetInt64("max_file_size"))
		files := c.NewFiles(repoService, fileService)

		// request context wrapper
		wrap := c.Wrapper(c.Config{
			UserFunc:     repoService.UserByRememberToken,
			Router:       r,
			ErrorHandler: errs.HandleError,
			Permissions:  permissions,
			Render:       renderer,
			Assets:       assets,
		})

		// routes
		r.NotFoundHandler = wrap(errs.NotFound)
		// TODO don't apply all middleware to static file server
		r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		r.Handle("/", wrap(pages.Home)).Methods("GET").Name("home")
		r.Handle("/auth/callback", wrap(auth.Callback)).Methods("GET")
		r.Handle("/logout", wrap(auth.Logout)).Methods("GET").Name("logout")
		r.Handle("/login", wrap(auth.Login)).Methods("GET").Name("login")
		r.Handle("/spaces", wrap(c.RequireUser, spaces.List)).Methods("GET").Name("spaces")
		r.Handle("/spaces/{spaceName}", wrap(c.RequireUser, spaces.Show)).Methods("GET").Name("space")
		r.Handle("/new-space", wrap(c.RequireAdmin, spaces.New)).Methods("GET").Name("new_space")
		r.Handle("/spaces", wrap(c.RequireAdmin, spaces.Create)).Methods("POST").Name("create_space")
		r.Handle("/spaces/{spaceName}/edit", wrap(c.RequireAdmin, spaces.Edit)).Methods("GET").Name("edit_space")
		r.Handle("/spaces/{spaceName}", wrap(c.RequireAdmin, spaces.Update)).Methods("PUT").Name("update_space")
		r.Handle("/spaces/{spaceName}/folders", wrap(c.RequireUser, spaces.CreateFolder)).Methods("POST").Name("create_folder")
		r.Handle("/folders/{folderID}", wrap(c.RequireUser, folders.Show)).Methods("GET").Name("folder")
		r.Handle("/folders/{folderID}/edit", wrap(c.RequireUser, folders.Edit)).Methods("GET").Name("edit_folder")
		r.Handle("/folders/{folderID}", wrap(c.RequireUser, folders.Update)).Methods("PUT").Name("update_folder")
		r.Handle("/folders/{folderID}/files", wrap(c.RequireUser, folders.UploadFile)).Methods("POST").Name("upload_file")
		r.Handle("/folders/{folderID}", wrap(c.RequireUser, folders.Delete)).Methods("DELETE").Name("delete_folder")
		r.Handle("/files/{fileID}", wrap(files.Download)).Methods("GET").Name("download_file")
		r.Handle("/files/{fileID}", wrap(c.RequireUser, files.Delete)).Methods("DELETE").Name("delete_file")
		r.Handle("/share/{folderID}:{folderSlug}", wrap(folders.Share)).Methods("GET").Name("share_folder")

		type WebSocketRequest struct {
			Type string
			Body struct {
				Route  string
				Params map[string]string
			}
		}

		type ClientData struct {
			User *models.User
		}

		wsRouter := &turborouter.Router[ClientData, map[string]string]{
			Deserializer: func(msg []byte) (string, map[string]string, error) {
				r := &WebSocketRequest{}
				if err := json.Unmarshal(msg, r); err != nil {
					return "", nil, err
				}
				return r.Body.Route, r.Body.Params, nil
			},
		}

		mu := sync.RWMutex{}
		var greetName string
		nameSwapper := func(str string) {
			mu.Lock()
			greetName = str
			mu.Unlock()
		}

		wsRouter.Add("home", func(c *turbo.Client[ClientData], params map[string]string) {
			nameSwapper(params["name"])
			c.Send(turbo.ReplaceMatch(
				".bc-avatar-text",
				`<span class="bc-avatar-text">Hi `+
					params["name"]+
					`(user: `+
					c.Data.User.Name+
					`)<span id="current-time"></span></span>`,
			))
		})

		hub := turbo.NewHub(turbo.Config[ClientData]{
			Responder: wsRouter,
		})

		r.Handle("/ws", wrap(func(ctx *ctx.Ctx) error {
			hub.Handle(ctx.Res, ctx.Req, func(c *turbo.Client[ClientData]) {
				c.Data = ClientData{User: ctx.User}
				c.Join(ctx.User.ID)
			})
			return nil
		}))

		tick := time.NewTicker(5 * time.Second)
		go func() {
			for {
				select {
				case <-tick.C:
					mu.RLock()
					name := greetName
					mu.RUnlock()
					hub.Broadcast(turbo.ReplaceMatch("h1.bc-toolbar-title", `<h1 class="bc-toolbar-title">Hi `+name+`</h1>`))
				}
			}
		}()

		// apply these before request reaches the router
		handler := middleware.Apply(r,
			middleware.Recover(func(err any) {
				if config.Production {
					logger.With(zap.Stack("stack")).Error(err)
				} else {
					logger.Error(err)
				}
			}),
			func(next http.Handler) http.Handler {
				return http.MaxBytesHandler(next, viper.GetInt64("max_file_size"))
			},
			// apply before ProxyHeaders to avoid invalid referer errors
			csrf.Protect(
				[]byte(config.CookieSecret),
				csrf.CookieName("deliver.csrf"),
				csrf.Path("/"),
				csrf.Secure(config.Production),
				csrf.SameSite(csrf.SameSiteStrictMode),
				csrf.FieldName("_csrf_token"),
			),
			middleware.MethodOverride(
				middleware.MethodFromHeader(middleware.MethodHeader),
				middleware.MethodFromForm(middleware.MethodParam),
			),
			middleware.If(config.Production, handlers.ProxyHeaders),
			middleware.SetRequestID(func() string {
				return ulid.Make().String()
			}),
			zaphttp.SetLogger(logger.Desugar()),
			zaphttp.LogRequests(logger.Desugar()),
			crumb.Enable(
				crumb.WithErrorHandler(func(err error) {
					logger.Error(err)
				}),
			),
		)

		// start server
		// TODO make timeouts configurable
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		server := graceful.WithDefaults(&http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Minute,
			WriteTimeout: 10 * time.Minute,
		})
		logger.Infof("starting server at %s", addr)
		if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
			logger.Fatal(err)
		}
		logger.Info("gracefully stopped server")
	},
}
