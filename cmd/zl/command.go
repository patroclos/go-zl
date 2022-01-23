package main

import (
	"github.com/spf13/cobra"
	"jensch.works/zl/cmd/zl/context"
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

	return cmd, ctx
}

func runRoot(cmd *cobra.Command, ctx *context.Context, args []string) error {
	return nil
}
