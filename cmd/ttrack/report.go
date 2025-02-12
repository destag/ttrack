package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/toggl"
)

var cmdReport = &cli.Command{
	Name:      "report",
	Usage:     "generate time report",
	ArgsUsage: "<project>",
	Action:    runReport,
}

func runReport(ctx *cli.Context) error {
	fmt.Println("Getting report")
	cfg := ctx.Context.Value(configKey).(*config.Config)
	c := toggl.NewClient(cfg.TogglToken.String())

	if ctx.NArg() != 1 {
		return cli.ShowSubcommandHelp(ctx)
	}

	project := ctx.Args().Get(0)
	if project == "" {
		return cli.Exit("project not provided", 1)
	}

	usr, err := c.GetUserInfo()
	if err != nil {
		return err
	}

	p, err := c.GetProject(usr.DefaultWorkspaceID, project)
	if err != nil {
		return err
	}

	tasks, err := c.GetSummaryReport(usr.DefaultWorkspaceID, p.ID)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", tasks)
	return nil
}
