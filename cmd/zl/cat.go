package main

import (
	"fmt"
	"strings"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
)

func makeCmdCat(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{}
	cmd.Use = "cat"
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

		fmt.Println(zl.Readme().String())
		return nil
	}
	return cmd
}

/*
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
*/
