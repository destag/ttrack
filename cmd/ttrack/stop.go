package main

import (
	"fmt"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/toggl"
	"github.com/urfave/cli/v2"
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

	fmt.Printf("Stopping tracking '%s'\n", te.Description)
	return c.StopTimeEntry(te)
}
