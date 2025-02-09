package main

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/github"
	"github.com/destag/ttrack/internal/toggl"
)

var cmdStart = &cli.Command{
	Name:      "start",
	Usage:     "start time tracker",
	ArgsUsage: "<project_name> <issue_id>",
	Action:    runStart,
}

func runStart(ctx *cli.Context) error {
	fmt.Println("Starting tracker")
	cfg := ctx.Context.Value(configKey).(*config.Config)
	c := toggl.NewClient(cfg.TogglToken.String())
	gh := github.NewClient(cfg.GithubToken.String())

	if ctx.NArg() != 2 {
		return cli.ShowSubcommandHelp(ctx)
	}

	project := ctx.Args().Get(0)
	if project == "" {
		return cli.Exit("project not provided", 1)
	}

	id := ctx.Args().Get(1)
	if id == "" {
		return cli.Exit("id not provided", 1)
	}

	togglProject, ok := cfg.Projects[project]
	if !ok {
		return cli.Exit("project not configured", 1)
	}

	issueID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	issue, err := gh.GetIssue(project, issueID)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("%s #%d", issue.Title, issue.Number)

	usr, err := c.GetUserInfo()
	if err != nil {
		return err
	}

	return c.StartTimeEntry(usr.DefaultWorkspaceID, title, togglProject)
}
