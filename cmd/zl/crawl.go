package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
	"jensch.works/zl/pkg/zettel/elemz"
	"jensch.works/zl/pkg/zettel/graph"
)

type depthCrawler struct {
	max int
}

func (spec depthCrawler) Crawl(c crawl.Node) crawl.RecurseMask {
	if spec.max > 0 && len(c.Path) > spec.max {
		return crawl.None
	}
	switch c.Reason.Refbox {
	case nil:
		fmt.Println(printZet(c.N.Z))
	default:
		rel := c.Reason.Refbox.Rel
		if strings.Contains(rel, " ") {
			rel = fmt.Sprintf("%#v", rel)
		}
		fmt.Print(zettel.MustFmt(c.N.Z, fmt.Sprintf("{{.Id}}  {{.Title}} parent:%s/%s/refbox[%s]\n", c.Path[0].N.Z.Id(), c.Reason.Mask.String(), c.Reason.Refbox.Rel)))
	}
	return crawl.All
}

func makeCmdCrawl(store zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "crawl"

	depth := cmd.Flags().IntP("depth", "d", 0, "max-depth to traverse to")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		g, err := graph.Make(store)
		if err != nil {
			return err
		}
		crawler := crawl.New(g, depthCrawler{max: *depth}.Crawl)
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

		scn := elemz.ListScanner(store)
		zets := make([]zettel.Z, 0, 16)
		for zet := range scn.Scan(os.Stdin) {
			zets = append(zets, zet)
		}
		crawler.Crawl(zets...)

		return nil
	}
	return cmd
}
