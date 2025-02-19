package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/github"
	"github.com/destag/ttrack/internal/toggl"
)

var cmdStart = &cli.Command{
	Name:            "start",
	Usage:           "start time tracker",
	ArgsUsage:       "<project_name> <issue_id>",
	Action:          runStart,
	HideHelpCommand: true,
}

func runStart(ctx *cli.Context) error {
	fmt.Println("Starting tracker")
	cfg := ctx.Context.Value(configKey).(*config.Config)
	c := toggl.NewClient(cfg.TogglToken.String())
	gh := github.NewClient(cfg.GithubToken.String())

	var project config.Project
	var id string
	for rgx, pr := range cfg.Projects {
		re := regexp.MustCompile(rgx)
		if matches := re.FindStringSubmatch(ctx.Args().Get(0)); len(matches) > 1 {
			project = pr
			id = matches[1]
			break
		}
	}

	if ctx.NArg() != 1 {
		return cli.ShowSubcommandHelp(ctx)
	}

	if ctx.Args().Get(0) == "" {
		return cli.Exit("project not provided", 1)
	}

	issueID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	issue, err := gh.GetIssue(project.Project, issueID)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("%s #%d", issue.Title, issue.Number)

	usr, err := c.GetUserInfo()
	if err != nil {
		return err
	}

	return c.StartTimeEntry(usr.DefaultWorkspaceID, title, project.Name)
}
