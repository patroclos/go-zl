package main

import (
	"fmt"
	"log"
	"strings"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

type cmdBacklinks struct {
	st zettel.Storage
}

func (c cmdBacklinks) Help() string {
	return ""
}

func (c cmdBacklinks) Synopsis() string {
	return "zettel"
}

func (c cmdBacklinks) Run(args []string) int {
	zets, err := c.st.Resolve(strings.Join(args, " "))
	if err != nil {
		log.Fatal(err)
	}
	zl, err := pickOne(zets)
	if err != nil {
		log.Fatal(err)
	}

	scn := scan.ListScanner(c.st)
	iter := c.st.Iter()
	for iter.Next() {
		zl2 := iter.Zet()
		for ref := range scn.Scan(strings.NewReader(zl2.Readme().Text)) {
			if ref.Id() == zl.Id() {
				fmt.Printf("%s  %s\n", zl2.Id(), zl2.Readme().Title)
			}
		}
	}

	return 0
}
