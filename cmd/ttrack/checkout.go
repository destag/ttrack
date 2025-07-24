package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/urfave/cli/v3"

	"github.com/destag/ttrack/internal/config"
	"github.com/destag/ttrack/internal/project"
	"github.com/destag/ttrack/internal/toggl"
)

var cmdCheckout = &cli.Command{
	Name:    "checkout",
	Aliases: []string{"co"},
	Usage:   "Checkout a branch for the current task",
	Action:  runCheckout,
}

func runCheckout(ctx context.Context, cmd *cli.Command) error {
	cfg := cmd.Root().Metadata["config"].(*config.Config)
	c := toggl.NewClient(cfg.TogglToken.String())
	te, err := c.GetCurrentTimeEntry()
	if err != nil {
		return err
	}

	if te.ID == 0 {
		return cli.Exit("No tracking in progress", 1)
	}

	proj, id, found := project.Find(cfg.Projects, te.Description)
	if !found {
		return cli.Exit("Could not find project for running task", 1)
	}

	if proj.BranchFormat == "" {
		return cli.Exit("Branch format not configured for this project", 1)
	}

	branchName := fmt.Sprintf(proj.BranchFormat, id)

	fmt.Printf("Checking out branch: %s\n", branchName)

	if err := exec.Command("git", "checkout", branchName).Run(); err != nil {
		fmt.Println("Branch does not exist, creating")
		return exec.Command("git", "checkout", "-b", branchName).Run()
	}

	return nil
}
