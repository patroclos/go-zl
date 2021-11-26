package main

import "github.com/spf13/cobra"

func makeCmdBacklinks() cobra.Command {
	cmd := cobra.Command{
		Use: "backlinks [zettelref]",
		Short: "Create a list of markdown refs of zettel containg ref to the argument",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
