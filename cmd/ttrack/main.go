package main

import (
	"context"
	"log"
	"os"

	"github.com/destag/ttrack/internal/config"
	"github.com/urfave/cli/v2"
)

type contextKey string

const (
	configKey contextKey = "config"
)

func main() {
	app := &cli.App{
		Name:                 "ttrack",
		Usage:                "track time in toggl",
		DefaultCommand:       "status",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:       "config",
				Value:      "~/.config/ttrack/config.yml",
				HasBeenSet: true,
				Action: func(ctx *cli.Context, s string) error {
					cfg, err := config.Load(s)
					if err != nil {
						return err
					}

					ctx.Context = context.WithValue(ctx.Context, configKey, cfg)
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			cmdStart,
			cmdStatus,
			cmdStop,
			cmdResume,
			cmdCompletion,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
