package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/mitchellh/cli"
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
		"": func() (cli.Command, error) {
			return cmdList{st: store}, nil
		},
		"list": func() (cli.Command, error) {
			return cmdList{st: store}, nil
		},
		"new": func() (cli.Command, error) {
			return cmdNew{st: store}, nil
		},
		"edit": func() (cli.Command, error) {
			return cmdEdit{st: store}, nil
		},
		"cat": func() (cli.Command, error) {
			return cmdCat{st: store}, nil
		},
		"blinks": func() (cli.Command, error) {
			return cmdBacklinks{st: store}, nil
		},
	}

	exit, err := c.Run()
	if err != nil {
		log.Println(err)
	}
	os.Exit(exit)
}
