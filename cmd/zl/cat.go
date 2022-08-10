package main

import (
	"fmt"
	"strings"

	"github.com/go-clix/cli"
	"git.jensch.dev/zl/pkg/zettel"
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
