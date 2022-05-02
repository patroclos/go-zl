package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-clix/cli"
	"gonum.org/v1/gonum/graph/community"
	"jensch.works/zl/pkg/zettel"
	zlg "jensch.works/zl/pkg/zettel/graph"
	"jensch.works/zl/pkg/zettel/scan"
)

func makeCmdSummary(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{
		Use:     "summary",
		Aliases: []string{"sum"},
	}
	showAll := cmd.Flags().BoolP("all", "a", false, "show all")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		if !isTerminal(os.Stdin) {
			return fmt.Errorf("not handling listings yet")
		}

		rels := make(map[string][2]int)

		g, idmap, errs := zlg.Make(st)
		for _, err := range errs {
			log.Println(err)
		}

		reduced := community.Modularize(g, 2, nil)
		fmt.Printf("Num Nodes: %d\n", g.Nodes().Len())
		fmt.Printf("Communities: %d\n\n", len(reduced.Communities()))

		comm := reduced.Communities()
		for i, com := range comm {
			if len(com) < 2 {
				continue
			}
			fmt.Printf("COMMUNITY %d:\n", i)
			for _, n := range com {
				fmt.Println(printZet(idmap[n.ID()]))
			}
			fmt.Println()
		}

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
			if !*showAll && num[0] == 1 {
				continue
			}
			fmt.Printf("%q relation has %d refs across %d boxes\n", rel, num[1], num[0])
		}

		return nil
	}
	return cmd
}
