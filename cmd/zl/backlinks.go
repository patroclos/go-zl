package main

import (
	"fmt"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
	"jensch.works/zl/pkg/zettel/graph"
)

func makeCmdBacklinks(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{}
	cmd.Use = "backlinks"
	cmd.Aliases = []string{"blinks"}
	cmd.Run = func(cmd *cli.Command, args []string) error {
		zets, err := st.Resolve(strings.Join(args, " "))
		if err != nil {
			return err
		}
		zl, err := pickOne(zets)
		if err != nil {
			return err
		}

		g, err := graph.Make(st)
		if err != nil {
			return err
		}
		crawl.New(g, func(n crawl.Node) crawl.RecurseMask {
			if len(n.Path) == 0 {
				return crawl.Inbound
			}
			fmt.Printf("%s  rel:%#v\n", n.N.Z, n.Reason.Refbox.Rel)
			return crawl.None
		}).Crawl(zl)

		return nil
	}
	return cmd
}
