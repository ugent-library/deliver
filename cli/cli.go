package cli

import (
	"github.com/caarlos0/env/v8"
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
	cobra.CheckErr(env.ParseWithOptions(&config, env.Options{
		Prefix: "DELIVER_",
	}))
	config.AfterLoad()
}

func initLogger() {
	if config.Production {
		l, err := zap.NewProduction()
		cobra.CheckErr(err)
		logger = l.Sugar()
	} else {
		l, err := zap.NewDevelopment()
		cobra.CheckErr(err)
		logger = l.Sugar()
	}
}

func Run() error {
	return rootCmd.Execute()
}
