package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/visibility"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
)

func makeCmdList(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{}
	cmd.Use = "list"
	cmd.Aliases = []string{"ls"}
	all := cmd.Flags().BoolP("all", "a", false, "disable taint filtering")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		isTerm := isTerminal(os.Stdin)

		printZ := func(n crawl.Node) crawl.RecurseMask {
			fmt.Println(printZet(n.Z))
			return crawl.None
		}
		view := visibility.TaintView(printZ, strings.Split(os.Getenv(`ZL_TOLERATE`), ","))

		var c crawl.Crawler
		if *all {
			c = crawl.New(st, printZ)
		} else {
			c = crawl.New(st, view)
		}
		if isTerm {
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
