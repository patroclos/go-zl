package main

import (
	"fmt"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
)

func makeCmdCrawl(store zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "crawl"

	depth := cmd.Flags().IntP("depth", "d", 0, "max-depth to traverse to")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		zets, err := store.Resolve(strings.Join(args, " "))
		if err != nil {
			return err
		}
		zet, err := pickOne(zets)
		if err != nil {
			return err
		}
		crawl.NewCrawler(store).Crawl(zet, func(c crawl.Crawl) crawl.RecurseMask {
			if *depth > 0 && len(c.Path) > *depth {
				return crawl.MaskNone
			}
			fmt.Println(zettel.MustFmt(c.Z, zettel.ListingFormat))
			return crawl.MaskAll
		})
		return nil
	}
	return cmd
}
