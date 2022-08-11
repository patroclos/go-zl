package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"git.jensch.dev/joshua/zl/cmd/zl/view"
	"git.jensch.dev/joshua/zl/pkg/visibility"
	"git.jensch.dev/joshua/zl/pkg/zettel"
	"git.jensch.dev/joshua/zl/pkg/zettel/crawl"
	"git.jensch.dev/joshua/zl/pkg/zettel/graph"
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

		listing, err := scanListing(bufio.NewScanner(os.Stdin), st)
		for _, zet := range listing {
			c.Crawl(zet)
		}
		return err
	}
	return cmd
}

func scanListing(scn *bufio.Scanner, st zettel.Resolver) ([]zettel.Z, error) {
	zets := []zettel.Z{}
	for scn.Scan() {
		zettel, err := st.Resolve(scn.Text())
		if err != nil {
			return zets, err
		}

		zets = append(zets, zettel...)
	}
	return zets, nil
}
