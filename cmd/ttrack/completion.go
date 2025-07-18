package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/destag/ttrack/internal/autocomplete"
)

var cmdCompletion = &cli.Command{
	Name:   "completion",
	Usage:  "Generate zsh autocompletion",
	Action: runCompletion,
}

func runCompletion(ctx context.Context, cmd *cli.Command) error {
	comp, err := autocomplete.EmbeddedFiles.ReadFile("zsh_autocomplete")
	if err != nil {
		return err
	}

	fmt.Println(string(comp))
	return nil
}
