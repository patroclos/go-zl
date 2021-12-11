package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/cmd/zl/view"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

func makeRootCommand(st storage.Storer) (*cobra.Command, *context.Context) {
	ctx := &context.Context{
		Template: zettel.ListPrettyStatusFormat,
		Store:    st,
	}

	labelspecs := make([]string, 0, 4)
	cmd := &cobra.Command{
		Use:   "zl",
		Short: "Personal Knowledge Jumpdrive",
		RunE: func(cmd *cobra.Command, args []string) error {
			specs := []zettel.Labelspec{}
			for _, ls := range labelspecs {
				spec, err := zettel.ParseLabelspec(ls)
				if err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), err)
					continue
				}

				specs = append(specs, *spec)
			}
			ctx.Labels = specs

			return runRoot(cmd, ctx, args)
		},
	}

	cmd.PersistentFlags().StringVarP(&ctx.Template, "template", "t", zettel.ListPrettyStatusFormat, "Customize zettellist output")
	cmd.PersistentFlags().StringSliceVarP(&labelspecs, "label", "l", nil, "Filter zettel against a labelspec")

	cmd.AddCommand(makeCmdNew())
	cmd.AddCommand(makeCmdMake())
	cmd.AddCommand(makeCmdBacklinks(ctx))
	cmd.AddCommand(view.MakeCommand(ctx))
	cmd.AddCommand(MakeGraphCommand(ctx))
	cmd.AddCommand(MakePromptCommand(ctx))

	return cmd, ctx
}

func isTerminal() bool {
	info, _ := os.Stdin.Stat()
	return (info.Mode() & os.ModeCharDevice) != 0
}

func labelFilter(ctx *context.Context, in <-chan zettel.Zettel) chan zettel.Zettel {
	ch := make(chan zettel.Zettel)
	go func() {
		defer close(ch)
		for x := range in {
			meta, err := x.Metadata()
			if err != nil {
				ch <- x
				continue
			}

			if zettel.RunSpecs(ctx.Labels, meta.Labels) {
				ch <- x
			}

		}
	}()
	return ch
}

func runRoot(cmd *cobra.Command, ctx *context.Context, args []string) error {
	var stream <-chan zettel.Zettel
	switch isTerminal() {
	case true:
		stream = storage.AllChan(ctx.Store)
	case false:
		scn := scan.ListScanner(ctx.Store)
		stream = scn.Scan(os.Stdin)
	}

	if len(ctx.Labels) > 0 {
		stream = labelFilter(ctx, stream)
	}

	for x := range stream {
		txt, err := zettel.Fmt(x, ctx.Template)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(txt)
	}
	return nil
}
