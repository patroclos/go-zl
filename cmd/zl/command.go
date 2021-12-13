package main

import (
	"github.com/spf13/cobra"
	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/cmd/zl/view"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			specs := []zettel.Labelspec{}
			for _, ls := range labelspecs {
				spec, err := zettel.ParseLabelspec(ls)
				if err != nil {
					return err
				}

				specs = append(specs, *spec)
			}
			ctx.Labels = specs
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoot(cmd, ctx, args)
		},
	}

	cmd.PersistentFlags().StringVarP(&ctx.Template, "template", "t", zettel.ListPrettyStatusFormat, "Customize zettellist output")
	cmd.PersistentFlags().StringSliceVarP(&labelspecs, "label", "l", nil, "Filter zettel against a labelspec")

	cmd.AddCommand(makeCmdNew())
	cmd.AddCommand(makeCmdMake())
	cmd.AddCommand(makeCmdList(ctx))
	cmd.AddCommand(makeCmdBacklinks(ctx))
	cmd.AddCommand(view.MakeCommand(ctx))
	cmd.AddCommand(MakeGraphCommand(ctx))
	cmd.AddCommand(MakePromptCommand(ctx))

	return cmd, ctx
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
	return runList(ctx, cmd, args)
}
