package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	configFile string
	config     Config
)

var logger *zap.SugaredLogger

var rootCmd = &cobra.Command{
	Use: "dilliver",
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")

	viper.SetEnvPrefix("dilliver")
	viper.SetDefault("s3.region", "us-east-1")
	viper.SetDefault("s3.bucket", "dilliver")
	viper.SetDefault("addr", "localhost:3002")
	viper.SetDefault("session.name", "dilliver")
	viper.SetDefault("session.max_age", 86400*30) // 30 days
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		cobra.CheckErr(viper.ReadInConfig())
	}

	viper.AutomaticEnv()

	cobra.CheckErr(viper.Unmarshal(&config))
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
