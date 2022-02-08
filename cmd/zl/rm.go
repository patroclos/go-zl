package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

func makeCmdRemove(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "remove"
	cmd.Aliases = []string{"rm"}
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

		if len(backlinks) > 0 {
			fmt.Println("Backlinks found:")
			for i := range backlinks {
				fmt.Printf("* %s  %s\n", backlinks[i].Id(), backlinks[i].Readme().Title)
			}
			fmt.Printf("Proceed anyway? y/N: ")

			var yn string
			_, err := fmt.Scanln(&yn)
			if err != nil || !strings.EqualFold(yn, "y") {
				return err
			}
		}

		if err := st.Remove(zet); err != nil {
			return err
		}
		return nil
	}
	return cmd
}
