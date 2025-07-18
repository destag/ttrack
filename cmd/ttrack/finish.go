package main

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/toggl"
)

var cmdFinish = &cli.Command{
	Name:    "finish",
	Aliases: []string{"f"},
	Usage:   "Stop current tracking",
	Action:  runStop,
}

func runStop(ctx context.Context, cmd *cli.Command) error {
	cfg := cmd.Root().Metadata[configKey].(*config.Config)
	c := toggl.NewClient(cfg.TogglToken.String())

	te, err := c.GetCurrentTimeEntry()
	if err != nil {
		return err
	}

	if te.ID == 0 {
		fmt.Println("No tracking to stop")
		return nil
	}

	fmt.Printf("Stopping tracking %s\n", color.GreenString(te.Description))
	return c.StopTimeEntry(te)
}
