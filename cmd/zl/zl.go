package main

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/spf13/cobra"

	"jensch.works/zl/pkg/graph"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/storage/filesystem"
	"jensch.works/zl/pkg/zettel"
)

func toTemplateData(zl zettel.Zettel) zettel.ZettelTemplate {
	tmpl := zettel.ZettelTemplate{
		Id:    string(zl.Id()),
		Title: zl.Title(),
	}

	return tmpl
}

func main() {
	rootCmd := cobra.Command{
		Use:   "zl",
		Short: "Personal Knowledge Manager",
	}

	var frmt string
	rootCmd.PersistentFlags().StringVarP(&frmt, "format", "f", "{{ .Id }}  {{ .Title }}", "zettel format string")

	cmdGraph := &cobra.Command{
		Use: "graph",
		Run: func(cmd *cobra.Command, args []string) {
			zlpath, ok := os.LookupEnv("ZLPATH")
			if !ok {
				panic("no ZLPATH environment variable set")
			}
			st := &filesystem.ZettelStorage{
				Directory: zlpath,
			}

			gv := graphviz.New()
			gv.SetLayout(graphviz.FDP)
			graph, err := graph.Plot(gv, st)
			if err != nil {
				log.Println(err)
				return
			}
			gv.RenderFilename(graph, graphviz.SVG, "test.svg")
		},
	}

	cmdList := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			zlpath, ok := os.LookupEnv("ZLPATH")
			if !ok {
				panic("no ZLPATH environment variable set")
			}
			st := filesystem.ZettelStorage{
				Directory: zlpath,
			}

			for _, zl := range storage.All(st) {
				data := toTemplateData(zl)
				txt, err := zettel.FormatZettel(data, frmt)
				if err != nil {
					log.Println(err)
					return
				}

				var inbox *string = nil
				for l,v := range zl.(*filesystem.Zettel).Meta.Labels {
					if l == "zl/inbox" {
						inbox = &v
						break
					}
				}

				if lnk := zl.(*filesystem.Zettel).Meta.Link; lnk != nil {
					fmt.Printf("LNK %s -> %s", lnk.A, lnk.B)
				}

				switch inbox {
				case nil:
					fmt.Println(txt)
				default:
					fmt.Println("📥", txt)
				}

				// fmt.Println(zl.(*filesystem.Zettel).Meta)
			}
		},
	}

	rootCmd.AddCommand(cmdList)
	rootCmd.AddCommand(cmdGraph)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return

	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		panic("no ZLPATH environment variable set")
	}
	st := filesystem.ZettelStorage{
		Directory: zlpath,
	}

	switch os.Args[1] {
	case "refs":
		for _, z := range storage.All(st) {
			txt, err := z.Text()
			if err != nil {
				continue
			}
			for _, r := range zettel.Refs(txt) {
				zr, err := st.Zettel(r)
				if err != nil {
					log.Printf("unable to resolve reference to %s", r)
					continue
				}
				fmt.Printf("%s => %s\n%s => %s\n\n", z.Id(), zr.Id(), z.Title(), zr.Title())
			}
		}
	default:
		for _, z := range storage.All(st) {
			fmt.Printf("%s  %s\n", z.Id(), z.Title())
		}
	}

}
