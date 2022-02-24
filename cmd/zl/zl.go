package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-clix/cli"
	"github.com/go-git/go-billy/v5/osfs"
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

	root := &cli.Command{}
	root.Use = "zl"
	root.Run = makeCmdList(store).Run
	root.AddCommand(makeCmdList(store))
	root.AddCommand(makeCmdCat(store))
	root.AddCommand(makeCmdNew(store))
	root.AddCommand(makeCmdEdit(store))
	root.AddCommand(makeCmdMetaEdit(store))
	root.AddCommand(makeCmdRemove(store))
	root.AddCommand(makeCmdBacklinks(store))
	root.AddCommand(makeCmdExport(store))
	root.AddCommand(makeCmdCrawl(store))
	root.AddCommand(makeCmdSearch(store))
	root.AddCommand(makeCmdLabel(store))

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
