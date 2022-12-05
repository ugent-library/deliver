package cmd

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
		services, err := models.NewServices(models.Config{
			DatabaseURL: viper.GetString("database_url"),
		})
		cobra.CheckErr(err)

		// router
		r := mux.NewRouter()
		r.StrictSlash(true)
		r.UseEncodedPath()
		r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

		// controllers
		pages := c.NewPages()
		spaces := c.NewSpaces(services.Repository)

		// request context wrapper
		wrap := c.Wrapper(r)

		// routes
		r.HandleFunc("/", wrap(pages.Home)).Methods("GET").Name("home")
		r.HandleFunc("/spaces", wrap(spaces.List)).Methods("GET").Name("spaces")
		r.HandleFunc("/spaces", wrap(spaces.Create)).Methods("POST").Name("create_space")

		// start server
		err = http.ListenAndServe("localhost:3002", r)
		cobra.CheckErr(err)
	},
}
