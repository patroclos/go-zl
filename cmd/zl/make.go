package main

import "github.com/spf13/cobra"

func makeCmdMake() *cobra.Command {
	cmd := &cobra.Command{
		Use: "make",
		Short: "create zettel and output ref or zettellist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
