package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"

	"github.com/destag/ttrack/internal/config"
)

type contextKey string

const (
	configKey contextKey = "config"
)

func main() {
	ver := "unknown"
	bi, ok := debug.ReadBuildInfo()
	if ok {
		ver = bi.Main.Version
	}

	app := &cli.App{
		Name:                 "ttrack",
		Usage:                "track time in toggl",
		DefaultCommand:       "status",
		EnableBashCompletion: true,
		Version:              ver,
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
