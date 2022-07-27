package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/go-clix/cli"
	"gonum.org/v1/gonum/graph/community"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/elemz"
	zlg "jensch.works/zl/pkg/zettel/graph"
)

func makeCmdSummary(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{
		Use:     "summary",
		Aliases: []string{"sum"},
		Short:   "Finds communities in the z-graph using the Louvain algorithm",
	}
	boxes := &cli.Command{
		Use:     "refbox",
		Aliases: []string{"rbox"},
		Short:   "Summarizes refbox relations (number of boxes and refs total per relation)",
		Run: func(cmd *cli.Command, args []string) error {
			rels := make(map[string]struct{ boxes, refs int })

			iter := st.Iter()
			if !isTerminal(os.Stdin) {
				listing, err := scanListing(bufio.NewScanner(os.Stdin), st)
				if err != nil {
					return fmt.Errorf("failed reading listing from stdin: %w", err)
				}
				iter = zettel.Slice(listing).Iter()
			}
			for iter.Next() {
				txt := iter.Zet().Readme().Text
				boxes := elemz.Refboxes(txt)
				for _, box := range boxes {
					x := rels[box.Rel]
					x.boxes++
					x.refs += len(box.Refs)
					rels[box.Rel] = x
				}
			}

			for rel, stats := range rels {
				if stats.boxes == 1 {
					continue
				}
				fmt.Printf("%q relation has %d refs across %d boxes\n", rel, stats.refs, stats.boxes)
			}
			return nil
		},
	}
	cmd.AddCommand(boxes)
	cmd.Run = func(cmd *cli.Command, args []string) error {
		g, err := zlg.Make(st)
		if err != nil {
			return fmt.Errorf("failed making graph: %w", err)
		}

		if !isTerminal(os.Stdin) {
			listing, err := scanListing(bufio.NewScanner(os.Stdin), st)
			if err != nil {
				return fmt.Errorf("failed reading listing from stdin: %w", err)
			}
			zs := zettel.Slice(listing)
			g, err = zlg.Make(zs)
			if err != nil {
				return fmt.Errorf("failed making small-graph: %w", err)
			}
			// * read listing from stdin
			// * prune all nodes not in the list from the graph
		}

		reduced := community.Modularize(g, 2, nil)
		fmt.Printf("Num Nodes: %d\n", len(g.Verts))
		fmt.Printf("Communities: %d\n\n", len(reduced.Communities()))

		comm := reduced.Communities()
		for i, com := range comm {
			if len(com) < 2 {
				continue
			}
			fmt.Printf("COMMUNITY %d:\n", i)
			for _, n := range com {
				fmt.Println(printZet(g.Verts[n.ID()].Z))
			}
			fmt.Println()
		}
		return nil
	}
	return cmd
}
