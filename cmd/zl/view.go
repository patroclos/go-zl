package main

import (
	"github.com/spf13/cobra"
	"jensch.works/zl/pkg/storage"
)

func makeCmdView(st storage.Storer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "view [zlref]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}
