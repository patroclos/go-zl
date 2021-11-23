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

	var wide bool
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
				if wide {
					frmt = zettel.DefaultWideFormat
				}
				txt, err := zettel.FormatZettel(zl, frmt)
				if err != nil {
					log.Println(err)
					return
				}

				fmt.Println(txt)
			}
		},
	}
	cmdList.Flags().BoolVarP(&wide, "wide", "w", false, "Use wide format")


	rootCmd.AddCommand(cmdList)
	rootCmd.AddCommand(cmdGraph)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
