package main

import (
	"fmt"
	"os"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

func makeCmdSummary(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{
		Use:     "summary",
		Aliases: []string{"sum"},
	}
	cmd.Run = func(cmd *cli.Command, args []string) error {
		if !isTerminal(os.Stdin) {
			return fmt.Errorf("not handling listings yet")
		}

		rels := make(map[string][2]int)

		iter := st.Iter()
		for iter.Next() {
			txt := iter.Zet().Readme().Text
			boxes := scan.All(txt)
			for _, box := range boxes {
				x := rels[box.Rel]
				x[0]++
				x[1] += len(box.Refs)
				rels[box.Rel] = x
			}
		}

		for rel, num := range rels {
			fmt.Printf("%q relation has %d refs across %d boxes\n", rel, num[1], num[0])
		}

		return nil
	}
	return cmd
}
