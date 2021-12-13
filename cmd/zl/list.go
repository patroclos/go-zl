package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
)

func makeCmdList(ctx *context.Context) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(ctx, cmd, args)
		},
	}
}

func runList(ctx *context.Context, cmd *cobra.Command, args []string) error {
	stream := storage.AllChan(ctx.Store)

	if len(ctx.Labels) > 0 {
		stream = labelFilter(ctx, stream)
	}

	for x := range stream {
		txt, err := zettel.Fmt(x, ctx.Template)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), txt)
	}
	return nil
}
