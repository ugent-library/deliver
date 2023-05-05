package cmd

import (
	"github.com/caarlos0/env/v8"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	// load .env file if present
	_ "github.com/joho/godotenv/autoload"

	// register objectstore backends
	_ "github.com/ugent-library/deliver/objectstore/s3"
)

var config Config

var logger *zap.SugaredLogger

var rootCmd = &cobra.Command{
	Use: "deliver",
}

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
	var l *zap.Logger
	var e error
	if config.Production {
		l, e = zap.NewProduction()
	} else {
		l, e = zap.NewDevelopment()
	}
	cobra.CheckErr(e)
	logger = l.Sugar()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
