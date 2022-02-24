package cmd

import "github.com/go-clix/cli"

func Command() *cli.Command {
	c := &cli.Command{
		Use:   "config",
		Short: "Read and manage zl configuration",
	}
	c.AddCommand(Switch())
	return c
}
