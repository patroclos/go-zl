package main

import (
	"fmt"

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
	iter := c.ctx.Store.Iter()
	for iter.Next() {
		zl := iter.Zet()
		fmt.Printf("%s  %s\n", zl.Id(), zl.Readme().Title)
	}
	return 0
}
