package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
	"jensch.works/zl/pkg/zettel/graph"
)

var uriRegex = regexp.MustCompile(`\w+:(\/?\/?)[^\s]+`)

func makeCmdBacklinks(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{}
	cmd.Use = "backlinks"
	cmd.Aliases = []string{"blinks"}
	cmd.Run = func(cmd *cli.Command, args []string) error {
		g, err := graph.Make(st)
		if err != nil {
			return err
		}
		q := strings.Join(args, " ")
		isUri := uriRegex.Match([]byte(q))
		if isUri {
			return fmt.Errorf("tbd")
		} // end-if isUri

		zets, err := st.Resolve(q)
		if err != nil {
			return err
		}
		zl, err := pickOne(zets)
		if err != nil {
			return err
		}

		idSelf := graph.Id(zl)
		inbound := g.To(idSelf)
		for inbound.Next() {
			idIn := inbound.Node().ID()
			rb := g.EdgeRefbox(idIn, idSelf)
			fmt.Printf("%s  rel:%#v\n", g.NodeZ(idIn).Z, rb.Rel)
		}

		crawl.New(g, func(n crawl.Node) crawl.RecurseMask {
			if len(n.Path) == 0 {
				return crawl.Inbound
			}
			return crawl.None
		}).Crawl(zl)

		return nil
	}
	return cmd
}
