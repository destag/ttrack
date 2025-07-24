package main

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/urfave/cli/v3"

	"github.com/destag/ttrack/internal/config"
)

var cmdConfig = &cli.Command{
	Name:   "config",
	Usage:  "Print current config",
	Action: runConfig,
}

func runConfig(ctx context.Context, cmd *cli.Command) error {
	cfg := cmd.Root().Metadata["config"].(*config.Config)
	out, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
