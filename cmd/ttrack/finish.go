package main

import (
	"context"
	"fmt"
	"time"

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

	dur := time.Since(te.Start).Truncate(time.Second)
	fmt.Printf(
		"%s stopping %s %s\n",
		stopIcon,
		color.WhiteString(te.Description),
		color.HiCyanString(dur.String()),
	)
	return c.StopTimeEntry(te)
}
