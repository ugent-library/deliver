package cli

import (
	// "github.com/caarlos0/env/v8"

	"os"

	"github.com/caarlos0/env/v8"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	// load .env file if present
	_ "github.com/joho/godotenv/autoload"

	// register objectstore backends
	_ "github.com/ugent-library/deliver/objectstore/s3"
)

var (
	config Config
	logger *zap.SugaredLogger
)

func initConfig() error {
	err := env.ParseWithOptions(&config, env.Options{
		Prefix: "DELIVER_",
	})
	if err != nil {
		return err
	}
	config.AfterLoad()
	return nil
}

func initLogger() error {
	var l *zap.Logger
	var err error
	if config.Production {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}

	if err != nil {
		return err
	}

	logger = l.Sugar()

	return nil
}

func Run() {
	app := &cli.App{
		Name:  "deliver",
		Usage: "Deliver CLI",
		Before: func(*cli.Context) error {
			if err := initConfig(); err != nil {
				return cli.Exit(err, 1)
			}
			if err := initLogger(); err != nil {
				return cli.Exit(err, 1)
			}
			return nil
		},
		Commands: []*cli.Command{
			foldersCmd,
			filesCmd,
			appCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
