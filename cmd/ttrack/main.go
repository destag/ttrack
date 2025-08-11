package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v3"

	"github.com/destag/ttrack/internal/config"
)

const (
	configKey = "config"
)

var debugMode bool

func main() {
	ver := "unknown"
	if info, ok := debug.ReadBuildInfo(); ok {
		ver = info.Main.Version
	}

	cmd := &cli.Command{
		Name:                  "ttrack",
		Usage:                 "track time in toggl",
		EnableShellCompletion: true,
		Version:               ver,
		DefaultCommand:        "status",
		Metadata:              make(map[string]any),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				HideDefault: true,
				Destination: &debugMode,
			},
			&cli.StringFlag{
				Name:  "config",
				Value: "~/.config/ttrack/config.yml",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			configPath := cmd.String("config")
			if debugMode {
				fmt.Printf("Loading config from %s\n", configPath)
			}
			cfg, err := config.Load(configPath)
			if err != nil {
				return ctx, err
			}
			cmd.Metadata[configKey] = cfg
			return ctx, nil
		},
		Commands: []*cli.Command{
			cmdStart,
			cmdStatus,
			cmdFinish,
			cmdResume,
			cmdCheckout,
			cmdConfig,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
