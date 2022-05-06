package main

import (
	"fmt"
	"os"

	"github.com/go-clix/cli"
	"gonum.org/v1/gonum/graph/community"
	"jensch.works/zl/pkg/zettel"
	"jensch.works/zl/pkg/zettel/graph"
)

func makeCmdGraph(st zettel.Storage) *cli.Command {
	cmd := new(cli.Command)
	cmd.Use = "graph"
	cmdHubs := new(cli.Command)
	cmdHubs.Use = "hubs"
	cmdHubs.Run = func(cmd *cli.Command, args []string) error {
		var domain zettel.ZettelerIter
		if !isTerminal(os.Stdin) {
			// TODO: set domain to stdin listing
		} else {
			domain = st
		}

		zg, nodemap, errs := graph.MakeG(st)
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}

		modular := community.Modularize(zg, 2, nil)

		// TODO: find hubs from communities
		_, _, _ = domain, nodemap, modular
		return nil
	}
	cmd.AddCommand(cmdHubs)

	return cmd
}
