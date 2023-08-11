package cli

import (
	"github.com/caarlos0/env/v8"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	// register objectstore backends
	_ "github.com/ugent-library/deliver/objectstore/s3"
)

var (
	config Config
	logger *zap.SugaredLogger

	rootCmd = &cobra.Command{
		Use:   "deliver",
		Short: "Deliver CLI",
	}
)

func init() {
	cobra.OnInitialize(initConfig, initLogger)
	cobra.OnFinalize(func() {
		logger.Sync()
	})
}

func initConfig() {
	cobra.CheckErr(env.Parse(&config))
}

func initLogger() {
	if config.Env == "local" {
		l, err := zap.NewDevelopment()
		cobra.CheckErr(err)
		logger = l.Sugar()
	} else {
		l, err := zap.NewProduction()
		cobra.CheckErr(err)
		logger = l.Sugar()
	}
}

func Run() error {
	return rootCmd.Execute()
}
