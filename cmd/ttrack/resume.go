package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/toggl"
)

var cmdResume = &cli.Command{
	Name:    "resume",
	Aliases: []string{"r"},
	Usage:   "Resumes last tracking",
	Action:  runResume,
}

func runResume(ctx context.Context, cmd *cli.Command) error {
	cfg := cmd.Root().Metadata[configKey].(*config.Config)
	c := toggl.NewClient(cfg.TogglToken.String())

	te, err := c.GetCurrentTimeEntry()
	if err != nil {
		return err
	}

	if te.ID != 0 {
		fmt.Println("Already tracking")
		return errors.New("tracking in progress")
	}

	tes, err := c.GetTimeEntries()
	if err != nil {
		return err
	}

	if len(tes) == 0 {
		fmt.Println("No time entries")
		return nil
	}

	te = tes[0]

	fmt.Println("Resuming tracker")
	return c.StartTimeEntry(te.WorkspaceID, te.Description, te.ProjectName)
}
