package main

import (
	"log"
	"os"

	"github.com/go-clix/cli"
)

func main() {
	root := cli.Command{
		Use: "zlmigrate",
	}
	root.AddCommand(makeCmdShortIds())
	root.AddCommand(makeCmdLint())

	if err := root.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
