package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

func makeCmdBacklinks(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backlinks [zettelref]",
		Short: "Create a list of markdown refs of zettel containg ref to the argument",
		RunE: func(cmd *cobra.Command, args []string) error {
			target, err := ctx.Store.Zettel(zettel.Id(args[0]))
			if err != nil {
				return fmt.Errorf("Zettel %v not found: %w", args[0], err)
			}

			for ref := range scan.Backrefs(target.Id(), ctx.Store) {
				fmt.Fprintln(cmd.OutOrStdout(), zettel.MustFmt(ref, ctx.Template))
			}
			return nil
		},
	}

	return cmd
}
