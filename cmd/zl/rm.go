package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/crawl"
)

func makeCmdRemove(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "remove"
	cmd.Aliases = []string{"rm"}
	frce := cmd.Flags().BoolP("force", "f", false, "skip confirmation and integrity cheks")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		q := strings.Join(args, " ")
		matches, err := st.Resolve(q)
		if err != nil {
			log.Fatal(err)
		}
		zet, err := pickOne(matches)
		if err != nil {
			log.Fatal(err)
		}

		if *frce {
			return st.Remove(zet)
		}

		backlinks := make([]zettel.Zettel, 0, 8)
		crawl.New(st, func(n crawl.Node) crawl.RecurseMask {
			if len(n.Path) == 0 {
				return crawl.Inbound
			}
			backlinks = append(backlinks, n.Z)
			return crawl.None
		}).Crawl(zet)

		listing := zettel.MustFmt(zet, zettel.ListingFormat)
		fmt.Fprintln(os.Stderr, listing)
		if len(backlinks) > 0 {
			fmt.Fprintf(os.Stderr, "Backlinks found:")
			for i := range backlinks {
				fmt.Fprintf(os.Stderr, "%s\n", zettel.MustFmt(backlinks[i], zettel.ListFormat))
			}
		}

		fmt.Fprintf(os.Stderr, "Really delete? y/N: ")

		var yn string
		_, err = fmt.Scanln(&yn)
		if err != nil || !strings.EqualFold(yn, "y") {
			return err
		}

		return st.Remove(zet)
	}
	return cmd
}
