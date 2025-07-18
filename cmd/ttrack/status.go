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

var cmdStatus = &cli.Command{
	Name:   "status",
	Usage:  "Check current tracking",
	Action: runStatus,
}

func runStatus(ctx context.Context, cmd *cli.Command) error {
	cfg := cmd.Root().Metadata[configKey].(*config.Config)
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
}
