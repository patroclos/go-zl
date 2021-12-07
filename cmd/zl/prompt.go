package main

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"jensch.works/zl/cmd/zl/context"
	"jensch.works/zl/pkg/prompt"
	"jensch.works/zl/pkg/storage"
)

func MakePromptCommand(ctx *context.Context) *cobra.Command {
	pc := promptCommand{ctx: ctx}
	cmd := &cobra.Command{
		Use:  "prompt CMD",
		RunE: pc.Run,
	}
	cmdList := &cobra.Command{
		Use:  "list",
		RunE: pc.List,
	}
	cmd.AddCommand(cmdList)
	return cmd
}

type promptCommand struct {
	ctx *context.Context
}

func (c promptCommand) Run(cmd *cobra.Command, args []string) error {
	// TODO: load all prompts, pick normal distribution + intention skew, display, review
	return nil
}

func (c promptCommand) List(cmd *cobra.Command, args []string) error {
	zets := storage.All(c.ctx.Store)
	wg := new(sync.WaitGroup)
	wg.Add(len(zets))
	for _, zet := range zets {
		zl := zet
		go func() {
			defer wg.Done()
			txt, err := zl.Text()
			if err != nil {
				return
			}
			prompts := prompt.ExtractAll(txt)
			for _, p := range prompts {
				fmt.Printf("%s\n\n", p.String())
			}
		}()
	}

	wg.Wait()

	return nil
}
