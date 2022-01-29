package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"jensch.works/zl/pkg/zettel"
)

type cmdList struct {
	st zettel.Storage
}

func (c cmdList) Help() string {
	return `Lists all zettel.
Using it as a filter makes it resolve every line of input and list every match.`
}

func (c cmdList) Synopsis() string {
	return ""
}

func (c cmdList) Run(args []string) int {
	isTerm := isTerminal(os.Stdin)
	if isTerm {
		iter := c.st.Iter()
		for iter.Next() {
			zl := iter.Zet()
			fmt.Printf("%s  %s\n", zl.Id(), zl.Readme().Title)
		}
		return 0
	}

	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		zets, err := c.st.Resolve(scn.Text())
		if err != nil {
			log.Println(err)
			continue
		}

		for _, zl := range zets {
			fmt.Printf("%s  %s\n", zl.Id(), zl.Readme().Title)
		}
	}
	if err := scn.Err(); err != nil {
		log.Println(err)
	}
	return 0
}
