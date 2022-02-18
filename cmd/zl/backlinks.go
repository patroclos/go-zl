package main

import (
	"fmt"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
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

		crawl.New(st, func(n crawl.Node) crawl.RecurseMask {
			if len(n.Path) == 0 {
				return crawl.Inbound
			}
			fmt.Println(n.Z)
			return crawl.None
		}).Crawl(zl)

		return nil
	}
	return cmd
}
