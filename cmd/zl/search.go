package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/storage/strutil"
	"jensch.works/zl/pkg/visibility"
	"jensch.works/zl/pkg/zettel"
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

		matches := 0
		for iter := st.Iter(); iter.Next(); {
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
				fmt.Println(zet)
			}
		}

		log.Printf("plain:%q labels:%v (%d matches)", plain.String(), labels, matches)
		return nil
	}
	return cmd
}
