package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/mitchellh/cli"
	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/pkg/storage"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		panic("no ZLPATH environment variable set")
	}

	dir := osfs.New(zlpath)
	store, err := storage.NewStore(dir)
	if err != nil {
		log.Fatal(err)
	}

	c := cli.NewCLI("zl", "0.1.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"new": func() (cli.Command, error) {
			return cmdNew{
				ctx: &context.Context{
					Store: store,
				},
			}, nil
		},
		"edit": func() (cli.Command, error) {
			return cmdEdit{
				ctx: &context.Context{
					Store: store,
				},
			}, nil
		},
	}

	exit, err := c.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exit)
}
