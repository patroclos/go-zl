package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"jensch.works/zl/cmd/zl/context"
)

type cmdList struct {
	ctx *context.Context
}

func (c cmdList) Help() string {
	return ""
}

func (c cmdList) Synopsis() string {
	return ""
}

func (c cmdList) Run(args []string) int {
	isTerm := isTerminal(os.Stdin)
	if isTerm {
		iter := c.ctx.Store.Iter()
		for iter.Next() {
			zl := iter.Zet()
			fmt.Printf("%s  %s\n", zl.Id(), zl.Readme().Title)
		}
		return 0
	}

	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		zets, err := c.ctx.Store.Resolve(scn.Text())
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
