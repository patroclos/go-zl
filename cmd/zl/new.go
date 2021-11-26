package main

import "github.com/spf13/cobra"

func makeCmdNew() cobra.Command {
	cmd := cobra.Command{
		Use: "new [title]",
		Short: "Create and edit new Zettel",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
