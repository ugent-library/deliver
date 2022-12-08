package cmd

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	c "github.com/ugent-library/dilliver/controllers"
	"github.com/ugent-library/dilliver/models"
)

func init() {
	rootCmd.AddCommand(appCmd)
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Start the web app server",
	Run: func(cmd *cobra.Command, args []string) {
		// setup services
		services, err := models.NewServices(models.Config{
			DatabaseURL: viper.GetString("database_url"),
		})
		if err != nil {
			logger.Fatal(err)
		}

		// setup router
		r := mux.NewRouter()
		r.StrictSlash(true)
		r.UseEncodedPath()
		r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

		// setup sessions
		sessionName := viper.GetString("session_name")
		sessionStore := sessions.NewCookieStore([]byte(viper.GetString("session_secret")))
		sessionStore.MaxAge(viper.GetInt("session_max_age"))
		sessionStore.Options.Path = "/"
		sessionStore.Options.HttpOnly = true
		sessionStore.Options.Secure = viper.GetBool("production")

		// controllers
		pages := c.NewPages()
		spaces := c.NewSpaces(services.Repository)

		// request context wrapper
		wrap := c.Wrapper(c.Config{
			Router:       r,
			SessionStore: sessionStore,
			SessionName:  sessionName,
		})

		// add routes
		r.HandleFunc("/", wrap(pages.Home)).Methods("GET").Name("home")
		r.HandleFunc("/spaces", wrap(spaces.List)).Methods("GET").Name("spaces")
		r.HandleFunc("/spaces", wrap(spaces.Create)).Methods("POST").Name("create_space")

		// start server
		if err = http.ListenAndServe(viper.GetString("app_addr"), r); err != nil {
			logger.Fatal(err)
		}
	},
}
