package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	zlpath, ok := os.LookupEnv("ZLPATH")
	if !ok {
		panic("no ZLPATH environment variable set")
	}
	fmt.Println(zlpath)

	/*
		rootCmd, ctx := makeRootCommand(st)
		// rootCmd.AddCommand(cmdGraph)
		// rootCmd.AddCommand(cmdPrompt)

		if err := rootCmd.Execute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	*/
}
