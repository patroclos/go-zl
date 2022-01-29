package main

import (
	"fmt"
	"log"
	"strings"

	"jensch.works/zl/pkg/zettel"
)

type cmdCat struct {
	st zettel.Storage
}

func (c cmdCat) Help() string {
	return `Renders the given zettel to the terminal`
}

func (c cmdCat) Synopsis() string {
	return "zettel"
}

func (c cmdCat) Run(args []string) int {
	q := strings.Join(args, " ")
	zets, err := c.st.Resolve(q)
	if err != nil {
		log.Fatal(err)
	}

	zl, err := pickOne(zets)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(zl.Readme().String())
	return 0
}
