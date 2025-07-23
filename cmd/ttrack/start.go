package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/github"
	"github.com/destag/ttrack/internal/jira"
	"github.com/destag/ttrack/internal/project"
	"github.com/destag/ttrack/internal/toggl"
)

var cmdStart = &cli.Command{
	Name:            "start",
	Aliases:         []string{"s"},
	Usage:           "Start time tracker",
	ArgsUsage:       "<project_name> <issue_id>",
	Action:          runStart,
	HideHelpCommand: true,
}

func runStart(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("Starting tracker")
	cfg := cmd.Root().Metadata[configKey].(*config.Config)
	c := toggl.NewClient(cfg.TogglToken.String())
	gh := github.NewClient(cfg.GithubToken.String())
	jc := jira.NewClient(cfg.Jira.Username, cfg.Jira.Token.String(), cfg.Jira.BaseURL)

	if cmd.NArg() != 1 {
		return cli.ShowAppHelp(cmd)
	}

	input := cmd.Args().Get(0)
	if input == "" {
		return cli.Exit("project not provided", 1)
	}

	proj, id, found := project.Find(cfg.Projects, input)
	if !found {
		return cli.Exit("project not found", 1)
	}

	var title string

	switch proj.Type {
	case "jira":
		task, err := jc.GetTask(id)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}
		title = fmt.Sprintf("%s %s", task.ID, task.Description)
	case "github":
		issueID, err := strconv.Atoi(id)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		issue, err := gh.GetIssue(proj.Project, issueID)
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}

		title = fmt.Sprintf("%s #%d", issue.Title, issue.Number)
	default:
		return cli.Exit("project type not supported", 1)
	}

	usr, err := c.GetUserInfo()
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	return c.StartTimeEntry(usr.DefaultWorkspaceID, title, proj.Name)
}
