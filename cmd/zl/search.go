package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"git.jensch.dev/joshua/zl/pkg/storage/strutil"
	"git.jensch.dev/joshua/zl/pkg/visibility"
	"git.jensch.dev/joshua/zl/pkg/zettel"
)

func makeCmdSearch(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "search"
	cmd.Args = cli.ArgsMin(1)
	cmd.Run = func(cmd *cli.Command, args []string) error {
		var plain strings.Builder

		sep := false
		labels := make([]zettel.Labelspec, 0)
		for i := range args {
			if strings.HasPrefix(args[i], "label:") {
				spec, err := zettel.ParseLabelspec(args[i][6:])
				if err != nil {
					return err
				}
				labels = append(labels, spec)
				continue
			}
			if sep {
				plain.WriteRune(' ')
			} else {
				sep = true
			}
			plain.WriteString(args[i])
		}

		iter := st.Iter()
		if !isTerminal(os.Stdin) {
			listing, err := scanListing(bufio.NewScanner(os.Stdin), st)
			if err != nil {
				return fmt.Errorf("failed reading listing from stdin: %w", err)
			}
			iter = zettel.Slice(listing).Iter()
		}

		matches := 0
		for iter.Next() {
			zet := iter.Zet()
			if !visibility.Visible(zet, strings.Split(os.Getenv("ZL_TOLERATE"), ",")) {
				continue
			}
			var veto bool
			for i := range labels {
				if !labels[i].Match(zet.Metadata().Labels) {
					veto = true
					break
				}
			}
			if veto {
				continue
			}

			if strutil.ContainsFold(zet.Readme().String(), plain.String()) {
				matches++
				fmt.Println(printZet(zet))
			}
		}

		log.Printf("plain:%q labels:%v (%d matches)", plain.String(), labels, matches)
		return nil
	}
	return cmd
}
