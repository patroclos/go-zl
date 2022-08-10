package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-clix/cli"
	"github.com/go-git/go-billy/v5/osfs"
	"git.jensch.dev/zl/pkg/prompt"
	"git.jensch.dev/zl/pkg/storage"
	"git.jensch.dev/zl/pkg/zettel/elemz"
)

func main() {
	elemz.DefaultParser.Parsers = append([]elemz.Parser{prompt.Parser()}, elemz.DefaultParser.Parsers...)

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

	ver := "dunno"
	info, ok := debug.ReadBuildInfo()
	if ok {
		ver = fmt.Sprintf("go:%q zl:%s", info.GoVersion, info.Main.Version)
	}
	root := &cli.Command{
		Use:     "zl",
		Version: ver,
		Run:     makeCmdList(store).Run,
	}
	format := root.Flags().StringP("format", "f", "listing", "Format string")
	root.Run = func(cmd *cli.Command, args []string) error {
		sub := makeCmdList(store)
		os.Args = append(os.Args[0:1], "-f", *format)
		return sub.Execute()
	}
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
	root.AddCommand(makeCmdTldr(store))
	root.AddCommand(makeCmdSummary(store))
	root.AddCommand(makeCmdPlace(store))
	root.AddCommand(makeCmdElem(store))

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
