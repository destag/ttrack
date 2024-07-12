package main

import (
	"fmt"

	"github.com/destag/ttrack/internal/autocomplete"
	"github.com/urfave/cli/v2"
)

var cmdCompletion = &cli.Command{
	Name:   "completion",
	Usage:  "generate zsh autocompletion",
	Action: runCompletion,
}

func runCompletion(ctx *cli.Context) error {
	comp, err := autocomplete.EmbeddedFiles.ReadFile("zsh_autocomplete")
	if err != nil {
		return err
	}

	fmt.Println(string(comp))
	return nil
}
