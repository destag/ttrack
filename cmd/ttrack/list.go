package main

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/jira"
	"github.com/destag/ttrack/internal/tasks"
)

var cmdList = &cli.Command{
	Name:            "list",
	Aliases:         []string{"l"},
	Usage:           "List available tasks",
	ArgsUsage:       "<project>",
	Action:          runList,
	HideHelpCommand: true,
}

func runList(ctx context.Context, cmd *cli.Command) error {
	cfg := cmd.Root().Metadata[configKey].(*config.Config)
	jc := jira.NewClient(cfg.Jira.Username, cfg.Jira.Token.String(), cfg.Jira.BaseURL)

	if cmd.NArg() != 1 {
		projs := slices.Collect(maps.Keys(cfg.Projects))
		return cli.Exit("project not specified, choose one of:\n"+strings.Join(projs, "\n"), 1)
	}

	proj, ok := cfg.Projects[cmd.Args().Get(0)]
	if !ok {
		return cli.Exit("project not found", 1)
	}

	var ts []*tasks.Task

	for _, task := range proj.Tasks {
		if task.Type != "jira" {
			continue
		}

		query := task.Query
		if query == "" {
			continue
		}

		out, err := jc.ListTasks(query)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		ts = append(ts, out...)
	}

	for _, task := range ts {
		fmt.Printf("%s %s\n", color.BlueString(task.ID), task.Description)
	}

	return nil
}
