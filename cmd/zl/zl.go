package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"jensch.works/zl/pkg/storage"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		panic("no ZLPATH environment variable set")
	}

	if len(os.Args) == 1 {
		log.Println("usage: zl SELECTOR")
		return
	}
	query := os.Args[1]
	fmt.Println(zlpath, query)

	dir := osfs.New(zlpath)
	store, err := storage.NewStore(dir)
	if err != nil {
		log.Fatal(err)
	}

	resolved, err := store.Resolve(query)
	if err != nil {
		log.Println(err)
	}
	log.Println(resolved)
}
