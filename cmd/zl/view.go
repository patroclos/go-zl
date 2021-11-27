package main

import (
	"github.com/spf13/cobra"
)

func makeCmdView(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "view [zlref]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}
