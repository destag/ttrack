package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/toggl"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type contextKey string

const (
	configKey contextKey = "config"
)

func main() {
	app := &cli.App{
		Name:           "ttrack",
		Usage:          "track time in toggl",
		DefaultCommand: "status",
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
			{
				Name:  "start",
				Usage: "start time tracker",
				Action: func(ctx *cli.Context) error {
					fmt.Println("starting tracker")
					fmt.Println("not implemented yet")
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "check current tracking",
				Action: func(ctx *cli.Context) error {
					cfg := ctx.Context.Value(configKey).(*config.Config)
					c := toggl.NewClient(cfg.TogglToken.String())

					te, err := c.GetCurrentTimeEntry()
					if err != nil {
						return err
					}

					if te.ID == 0 {
						fmt.Println("No tracking")
					} else {
						dur := time.Since(te.Start).Truncate(time.Second)
						fmt.Printf("Tracking %s %s\n",
							color.GreenString(te.Description),
							color.WhiteString(dur.String()))
					}

					return nil
				},
			},
			{
				Name:  "stop",
				Usage: "stop current tracking",
				Action: func(ctx *cli.Context) error {
					cfg := ctx.Context.Value(configKey).(*config.Config)
					c := toggl.NewClient(cfg.TogglToken.String())

					te, err := c.GetCurrentTimeEntry()
					if err != nil {
						return err
					}

					if te.ID == 0 {
						fmt.Println("No tracking to stop")
						return nil
					}

					fmt.Printf("Stopping tracking '%s'\n", te.Description)
					return c.StopTimeEntry(te)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
