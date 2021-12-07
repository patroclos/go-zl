package main

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/spf13/cobra"
	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/pkg/graph/gviz"
	"jensch.works/zl/pkg/storage"
	"jensch.works/zl/pkg/zettel/scan"
)

func MakeGraphCommand(ctx *context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "graph [listfile]",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cz := storage.AllChan(ctx.Store)
		if len(args) > 0 {
			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()
			cz = scan.ListScanner(ctx.Store).Scan(file)
		}

		gv := graphviz.New()
		gv.SetLayout(graphviz.DOT)

		net, err := gviz.Plot(gv, ctx.Store)
		if err != nil {
			return err
		}

		nodes := make(map[string]*cgraph.Node)
		for zl := range cz {
			id := string(zl.Id())
			node, err := net.CreateNode(id)
			if err != nil {
				log.Println(err)
				continue
			}
			node.SetLabel(zl.Title())
			node.SetURL("")
			nodes[id] = node
		}

		for zl := range storage.AllChan(ctx.Store) {
			txt, err := zl.Text()
			if err != nil {
				continue
			}
			refs := scan.Refs(txt)
			if len(refs) == 0 {
				id := string(zl.Id())
				node := nodes[id]
				delete(nodes, id)
				net.DeleteNode(node)
				continue
			}
			for _, ref := range refs {
				zref, err := ctx.Store.Zettel(ref)
				if err != nil {
					continue
				}
				id, refId := string(zl.Id()), string(zref.Id())
				name := fmt.Sprintf("%s-%s", string(id), string(refId))
				rnode, ok := nodes[refId]
				if !ok {
					continue
				}
				edge, err := net.CreateEdge(name, nodes[id], rnode)
				if err != nil {
					panic("possible?")
				}

				txt1, err1 := zl.Text()
				txt2, err2 := zref.Text()

				if err1 != nil || err2 != nil {
					continue
				}

				var weight float64
				if len(txt2) == 0 {
					weight = 1
				} else {
					weight = float64(len(txt1)) / float64(len(txt2))
				}
				log.Printf("[weight %v -> %v]: %04f", zl.Id(), zref.Id(), weight)
				edge.SetWeight(weight)
			}
		}

		switch len(args) {
		case 0:
			gv.RenderFilename(net, graphviz.PNG, "knet.png")
		case 1:
			gv.RenderFilename(net, graphviz.PNG, args[0])
		default:
			return fmt.Errorf("too many arguments, see usage")
		}

		return nil
	}
	return cmd
}
