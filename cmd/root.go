package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	// register objectstore backends
	_ "github.com/ugent-library/deliver/objectstore/s3"
)

var (
	configFile string
	config     Config
)

var logger *zap.SugaredLogger

var rootCmd = &cobra.Command{
	Use: "deliver",
}

func init() {
	viper.SetEnvPrefix("deliver")
	viper.SetDefault("port", 3000)
	viper.SetDefault("storage.backend", "s3")
	viper.SetDefault("max_file_size", 2_000_000_000)

	cobra.OnInitialize(initConfig, initLogger)
	cobra.OnFinalize(func() {
		logger.Sync()
	})

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		cobra.CheckErr(viper.ReadInConfig())
	}

	viper.AutomaticEnv()

	cobra.CheckErr(viper.Unmarshal(&config))

	if !config.Production && config.Banner == "" {
		config.Banner = "development"
	}
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
