package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/elemz"
)

// `elem [-t elem-type[,elem-type...]] <zref>` prints a summary of the text-elements
// contained in the referenced z(s)

func makeCmdElem(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "elements"
	cmd.Aliases = []string{"elem"}

	types := cmd.Flags().StringArrayP("type", "t", nil, "Set of element types to include")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		// listing in first
		var zets []zettel.Z
		if isTerminal(os.Stdin) {
			zets = zettel.All(st)
		} else {
			zs, err := scanListing(bufio.NewScanner(os.Stdin), st)
			if err != nil {
				return fmt.Errorf("reading z listing failed: %w", err)
			}
			zets = zs
		}

		for _, z := range zets {
			elems, err := elemz.Read(z.Readme().Text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed reading elements of %s\n", z.Id())
				continue
			}

			fmt.Printf("%s\n", lipgloss.NewStyle().Bold(true).Render(zettel.MustFmt(z, zettel.ListingFormat)))
			for i, el := range elems {
				ok := len(*types) == 0
				for _, t := range *types {
					if el.ElemType() == elemz.ElemType(t) {
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
				fmt.Printf("[%d] %s %v\n%s\n\n", i, el.ElemType(), el.Span(), el)
			}
			fmt.Println()
		}
		return nil
	}
	return cmd
}
