package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/go-clix/cli"
	"jensch.works/zl/pkg/zettel"
)

func makeCmdList(st zettel.Storage) *cli.Command {
	cmd := &cli.Command{}
	cmd.Use = "list"
	cmd.Aliases = []string{"ls"}
	cmd.Run = func(cmd *cli.Command, args []string) error {
		isTerm := isTerminal(os.Stdin)
		if isTerm {
			iter := st.Iter()
			for iter.Next() {
				zl := iter.Zet()
				fmt.Printf("%s  %s\n", zl.Id(), zl.Readme().Title)
			}
			return nil
		}

		scn := bufio.NewScanner(os.Stdin)
		for scn.Scan() {
			zets, err := st.Resolve(scn.Text())
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
		return nil
	}
	return cmd
}
