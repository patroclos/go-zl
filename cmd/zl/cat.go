package main

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

func makeCmdCat(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{}
	cmd.Use = "cat"
	scn := cmd.Flags().BoolP("scan", "s", false, "")
	cmd.Run = func(cmd *cli.Command, args []string) error {
		q := strings.Join(args, " ")
		zets, err := st.Resolve(q)
		if err != nil {
			return err
		}

		zl, err := pickOne(zets)
		if err != nil {
			return err
		}

		if *scn {
			elems, err := scan.Elements(st, zl.Readme().Text)
			if err != nil {
				log.Println(err)
			}
			lines := make([]string, 0)
			for scn := bufio.NewScanner(strings.NewReader(zl.Readme().Text)); scn.Scan(); {
				lines = append(lines, scn.Text())
			}
			for i := range elems {
				span := elems[i].Span
				log.Printf("-- ELEMENT[%d] (type %v) --\n", i, elems[i].Type)
				log.Println(strings.Join(lines[span.Start:span.Pos], "\n"))
			}
			return nil
		}

		fmt.Println(zl.Readme().String())
		return nil
	}
	return cmd
}
