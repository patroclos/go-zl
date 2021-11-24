package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/goccy/go-graphviz"
	"github.com/spf13/cobra"

	"jensch.works/zl/pkg/graph"
	"jensch.works/zl/pkg/prompt"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/storage/filesystem"
	"jensch.works/zl/pkg/zettel"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		panic("no ZLPATH environment variable set")
	}
	st := &filesystem.ZettelStorage{
		Directory: zlpath,
	}

	rootCmd := cobra.Command{
		Use:   "zl",
		Short: "Personal Knowledge Manager",
	}

	var frmt string
	rootCmd.PersistentFlags().StringVarP(&frmt, "format", "f", "{{ .Id }}  {{ .Title }}", "zettel format string")

	cmdGraph := &cobra.Command{
		Use: "graph",
		Run: func(cmd *cobra.Command, args []string) {

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

	cmdEdit := &cobra.Command{
		Use: "edit",
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")
			for _, zl := range storage.All(st) {
				if string(zl.Id()) == query {
					temp, err := os.CreateTemp("", "tmpzl*.md")
					if err != nil {
						return err
					}
					defer os.Remove(temp.Name())
					txt, err := zl.Text()
					if err != nil {
						return err
					}
					fmt.Fprintf(temp, "# %s\n\n%s", zl.Title(), txt)
					if err != nil {
						return err
					}
					cmd := exec.Command("vim", temp.Name())
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					err = cmd.Run()
					if exitErr, ok := err.(*exec.ExitError); ok {
						log.Println(exitErr, "aborting...")
						os.Exit(1)
					} else {
						if err != nil {
							return err
						}
					}

					temp.Seek(0, io.SeekStart)
					bytes, err := ioutil.ReadAll(temp)
					if err != nil {
						return err
					}

					zl.SetText(string(bytes))
					err = st.SetZettel(zl)
					if err != nil {
						return err
					}
					return nil
				}
			}

			return errors.New("not found")
		},
	}

	cmdPrompt := &cobra.Command{
		Use: "prompt",
		RunE: func(cmd *cobra.Command, args []string) error {
			allPrompts := make([]prompt.EmbeddedPrompt, 0, 256)
			err := st.ForEach(func(z zettel.Zettel) error {
				txt, err := z.Text()
				if err != nil {
					return nil
				}

				prompts := prompt.ExtractAll(txt)

				if len(prompts) == 0 {
					return nil
				}

				allPrompts = append(allPrompts, prompts...)

				return nil
			})
			if err != nil {
				return err
			}
			p := allPrompts[rand.Intn(len(allPrompts))]
			fmt.Printf("Q. %s\nA. ", p.Q)
			reader := bufio.NewReader(os.Stdin)
			_, err = reader.ReadString('\n')

			if err != nil {
				return err
			}

			fmt.Printf("A. %s\n", p.A)
			return nil
		},
	}

	rootCmd.AddCommand(cmdList)
	rootCmd.AddCommand(cmdEdit)
	rootCmd.AddCommand(cmdGraph)
	rootCmd.AddCommand(cmdPrompt)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
