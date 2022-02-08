package main

import (
	"fmt"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
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

		scn := scan.ListScanner(st)
		for iter := st.Iter(); iter.Next(); {
			zl2 := iter.Zet()
			for ref := range scn.Scan(strings.NewReader(zl2.Readme().Text)) {
				if ref.Id() == zl.Id() {
					fmt.Printf("%s  %s\n", zl2.Id(), zl2.Readme().Title)
				}
			}
		}

		return nil
	}
	return cmd
}
