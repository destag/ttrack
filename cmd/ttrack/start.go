package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/urfave/cli/v3"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/github"
	"github.com/destag/ttrack/internal/jira"
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

	var project config.Project
	var id string
	for rgx, pr := range cfg.Projects {
		re := regexp.MustCompile(rgx)
		if matches := re.FindStringSubmatch(cmd.Args().Get(0)); len(matches) > 1 {
			project = pr
			id = matches[1]
			break
		}
	}

	if cmd.NArg() != 1 {
		return cli.ShowAppHelp(cmd)
	}

	if cmd.Args().Get(0) == "" {
		return cli.Exit("project not provided", 1)
	}

	var title string

	switch project.Type {
	case "jira":
		task, err := jc.GetTask(cmd.Args().Get(0))
		if err != nil {
			return err
		}
		title = fmt.Sprintf("%s %s", task.ID, task.Description)
	case "github":
		issueID, err := strconv.Atoi(id)
		if err != nil {
			return err
		}

		issue, err := gh.GetIssue(project.Project, issueID)
		if err != nil {
			return err
		}

		title = fmt.Sprintf("%s #%d", issue.Title, issue.Number)
	default:
		return cli.Exit("project type not supported", 1)
	}

	usr, err := c.GetUserInfo()
	if err != nil {
		return err
	}

	return c.StartTimeEntry(usr.DefaultWorkspaceID, title, project.Name)
}
