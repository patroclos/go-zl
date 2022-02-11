package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
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

		backlinks := make([]zettel.Zettel, 0, 8)
		scn := scan.ListScanner(st)
		for iter := st.Iter(); iter.Next(); {
			other := iter.Zet()

			for ref := range scn.Scan(strings.NewReader(other.Readme().Text)) {
				if ref.Id() == zet.Id() {
					backlinks = append(backlinks, other)
					break
				}
			}
		}

		if *frce {
			return st.Remove(zet)
		}

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
