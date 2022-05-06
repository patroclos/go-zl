package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/cmd/zl/view"
	"jensch.works/zl/pkg/visibility"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
	"jensch.works/zl/pkg/zettel/graph"
)

func makeCmdList(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{}
	cmd.Use = "list"
	cmd.Aliases = []string{"ls"}
	all := cmd.Flags().BoolP("all", "a", false, "disable taint filtering")
	format := cmd.Flags().StringP("format", "f", "listing", "listing zettel format")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		termInput := isTerminal(os.Stdin)

		g, err := graph.Make(st)
		if err != nil {
			return err
		}

		i := 0
		printZ := func(n crawl.Node) crawl.RecurseMask {
			i++
			v := &view.Listing{
				Zets:  []zettel.Z{n.N.Z},
				Fmt:   *format,
				Dest:  os.Stdout,
				Color: true,
			}
			err := v.Render()
			if err != nil {
				// fmt.Fprintf(os.Stderr, "error listing %s: %v", n.N.Z.Id(), err)
				panic(*format)
			}
			return crawl.None
		}
		view := visibility.TaintView(printZ, strings.Split(os.Getenv(`ZL_TOLERATE`), ","))

		var c crawl.Crawler
		if *all {
			c = crawl.New(g, printZ)
		} else {
			c = crawl.New(g, view)
		}
		if termInput {
			for iter := st.Iter(); iter.Next(); {
				c.Crawl(iter.Zet())
			}
			return nil
		}

		scn := bufio.NewScanner(os.Stdin)
		for scn.Scan() {
			zets, err := st.Resolve(scn.Text())
			if err != nil {
				log.Println(err)
				continue
			}

			for _, zl := range zets {
				c.Crawl(zl)
			}
		}
		if err := scn.Err(); err != nil {
			log.Println(err)
		}
		return nil
	}
	return cmd
}
