package view

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel"
)

func MakeCommand(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "view [zlref]",
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")

			matches := make([]zettel.Zettel, 0, 8)

			for zl := range storage.AllChan(ctx.Store) {
				if meta, err := zl.Metadata(); err == nil {
					accept := zettel.RunSpecs(ctx.Labels, meta.Labels)
					if !accept {
						continue
					}
				}

				surface, err := zettel.Fmt(zl, "{{.Id}}  {{.Title}}")
				if err != nil {
					log.Println(err)
					continue
				}
				if strings.Contains(surface, query) {
					matches = append(matches, zl)
				}
			}

			l := len(matches)
			log.Printf("Found %d matches\n", l)
			switch l {
			case 0:
				os.Exit(1)
			case 1:
				viewZettel(matches[0])
			default:
				picked, err := pickOne(matches)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

				viewZettel(picked)
			}
			if l == 0 {
				os.Exit(1)
			}

			return nil
		},
	}
	return cmd
}
