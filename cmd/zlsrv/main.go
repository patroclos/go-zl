package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"jensch.works/zl/pkg/storage"
)

type UrlMaker struct {
	Base *url.URL
}

func (x UrlMaker) MakeUrl(feed []string, focus *string) *url.URL {
	if focus == nil {
		url, _ := x.Base.Parse(fmt.Sprintf("%s", strings.Join(feed, ",")))
		return url
	}

	url, _ := x.Base.Parse(fmt.Sprintf("%s#%s", strings.Join(feed, ","), *focus))
	return url
}

func main() {
	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		log.Fatal("ZLPATH evnironment variable is mandatory.")
	}

	dir := osfs.New(zlpath)
	store, err := storage.NewStore(dir)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := NewServer(store)
	if err != nil {
		log.Fatal(err)
	}

	addr := ":8000"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	if err := srv.Run(addr); err != nil {
		log.Fatal(err)
	}
}
