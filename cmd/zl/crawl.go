package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
	"jensch.works/zl/pkg/zettel/scan"
)

type depthCrawler struct {
	max int
}

func (spec depthCrawler) Crawl(c crawl.Node) crawl.RecurseMask {
	if spec.max > 0 && len(c.Path) > spec.max {
		return crawl.None
	}
	fmt.Println(printZet(c.Z))
	return crawl.All
}

func makeCmdCrawl(store zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "crawl"

	depth := cmd.Flags().IntP("depth", "d", 0, "max-depth to traverse to")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		crawler := crawl.New(store, depthCrawler{max: *depth}.Crawl)
		if isTerminal(os.Stdin) {
			zets, err := store.Resolve(strings.Join(args, " "))
			if err != nil {
				return err
			}
			zet, err := pickOne(zets)
			if err != nil {
				return err
			}

			crawler.Crawl(zet)
			return nil
		}

		scn := scan.ListScanner(store)
		zets := make([]zettel.Z, 0, 16)
		for zet := range scn.Scan(os.Stdin) {
			zets = append(zets, zet)
		}
		crawler.Crawl(zets...)

		return nil
	}
	return cmd
}
