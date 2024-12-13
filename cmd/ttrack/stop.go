package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/toggl"
)

var cmdStop = &cli.Command{
	Name:   "stop",
	Usage:  "stop current tracking",
	Action: runStop,
}

func runStop(ctx *cli.Context) error {
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

	fmt.Printf("Stopping tracking %s\n", color.GreenString(te.Description))
	return c.StopTimeEntry(te)
}
