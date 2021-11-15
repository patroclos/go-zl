package main

import (
	"fmt"
	"os"

	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/storage/filesystem"
)

func main() {
	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		panic("no ZLPATH environment variable set")
	}
	st := filesystem.ZettelStorage{
		Directory: zlpath,
	}

	for _, z := range storage.All(st) {
		fmt.Printf("%s  %s\n", z.Id(), z.Title())
	}
}
