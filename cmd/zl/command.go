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

func runRoot(cmd *cobra.Command, ctx *Context, args []string) error {
	var base <-chan zettel.Zettel
	switch isTerminal() {
	case true:
		base = storage.AllChan(ctx.st)
	case false:
		scn := scan.ListScanner(ctx.st)
		base = scn.Scan(os.Stdin)
	}

	if len(ctx.labels) > 0 {
		b2 := make(chan zettel.Zettel)
		old := base
		go func(){
			defer close(b2)
			for x := range old {
				meta, err := x.Metadata()
				if err != nil {
					b2 <- x
					continue
				}
				if zettel.RunSpecs(ctx.labels, meta.Labels) {
					b2 <- x
				}
			}
		}()

		base = b2
	}

	for x := range base {
		fmt.Printf("* %s  %s\n", x.Id(), x.Title())
	}
	return nil
}
