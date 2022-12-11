package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var configFile string

var logger *zap.SugaredLogger

var rootCmd = &cobra.Command{
	Use: "dilliver",
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")

	viper.SetEnvPrefix("dilliver")
	viper.SetDefault("app_addr", "localhost:3002")
	viper.SetDefault("session_name", "dilliver")
	viper.SetDefault("session_max_age", 86400*30) // 30 days
	viper.SetDefault("s3_bucket", "dilliver")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		cobra.CheckErr(viper.ReadInConfig())
	}

	viper.AutomaticEnv()
}

func initLogger() {
	var l *zap.Logger
	var e error
	if viper.GetBool("production") {
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
