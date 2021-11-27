package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

var (
	defaultFormat = "* {{.Id}}  {{.Title}}"
)

type Context struct {
	template string
	st       storage.Storer
	labels   []zettel.Labelspec
}

func makeRootCommand(st storage.Storer) (*cobra.Command, *Context) {
	ctx := &Context{
		template: defaultFormat,
		st:       st,
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
			ctx.labels = specs

			return runRoot(cmd, ctx, args)
		},
	}

	cmd.PersistentFlags().StringVarP(&ctx.template, "template", "t", defaultFormat, "Customize zettellist output")
	cmd.Flags().StringSliceVarP(&labelspecs, "label", "l", nil, "Filter zettel against a labelspec")

	cmd.AddCommand(makeCmdNew())
	cmd.AddCommand(makeCmdMake())
	cmd.AddCommand(makeCmdBacklinks())
	cmd.AddCommand(makeCmdView(ctx))

	return cmd, ctx
}

func isTerminal() bool {
	info, _ := os.Stdin.Stat()
	return (info.Mode() & os.ModeCharDevice) != 0
}

func labelFilter(ctx *Context, in <-chan zettel.Zettel) chan zettel.Zettel {
	ch := make(chan zettel.Zettel)
	go func() {
		defer close(ch)
		for x := range in {
			meta, err := x.Metadata()
			if err != nil {
				ch <- x
				continue
			}

			if zettel.RunSpecs(ctx.labels, meta.Labels) {
				ch <- x
			}

		}
	}()
	return ch
}

func runRoot(cmd *cobra.Command, ctx *Context, args []string) error {
	var stream <-chan zettel.Zettel
	switch isTerminal() {
	case true:
		stream = storage.AllChan(ctx.st)
	case false:
		scn := scan.ListScanner(ctx.st)
		stream = scn.Scan(os.Stdin)
	}

	if len(ctx.labels) > 0 {
		stream = labelFilter(ctx, stream)
	}

	for x := range stream {
		fmt.Printf("* %s  %s\n", x.Id(), x.Title())
	}
	return nil
}
