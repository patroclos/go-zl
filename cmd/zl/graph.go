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

		g, err := graph.Make(st)
		if err != nil {
			panic(err)
		}

		modular := community.Modularize(g, 2, nil)

		// TODO: find hubs from communities
		_, _ = domain, modular
		return fmt.Errorf("tbd")
	}
	cmd.AddCommand(cmdHubs)

	return cmd
}
