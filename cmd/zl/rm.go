package main

import (
	"fmt"
	"log"
	"strings"

	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/scan"
)

type cmdRemove struct {
	st zettel.Storage
}

func (rm cmdRemove) Help() string {
	return ``
}

func (rm cmdRemove) Synopsis() string {
	return ``
}

func (rm cmdRemove) Run(args []string) int {
	q := strings.Join(args, " ")
	matches, err := rm.st.Resolve(q)
	if err != nil {
		log.Fatal(err)
	}
	zet, err := pickOne(matches)
	if err != nil {
		log.Fatal(err)
	}

	backlinks := make([]zettel.Zettel, 0, 8)
	scn := scan.ListScanner(rm.st)
	for iter := rm.st.Iter(); iter.Next(); {
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
			return 1
		}
	}

	if err := rm.st.Remove(zet); err != nil {
		log.Fatal(err)
	}
	return 0
}
