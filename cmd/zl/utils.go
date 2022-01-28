package main

import (
	"fmt"
	"os"

	"jensch.works/zl/pkg/zettel"
)

func pickOne(zets []zettel.Zettel) (zettel.Zettel, error) {
	switch len(zets) {
	case 0:
		return nil, fmt.Errorf("no zettels to pick")
	case 1:
		return zets[0], nil
	}
	for i, z := range zets {
		fmt.Printf("[%d]: %s  %s\n", i+1, z.Id(), z.Readme().Title)
	}

	var idx int
	_, err := fmt.Scanln(&idx)
	if err != nil {
		return nil, err
	}

	if idx--; idx < 0 || idx >= len(zets) {
		return nil, fmt.Errorf("invalid index")
	}
	return zets[idx], nil
}

func isTerminal(f *os.File) bool {
	o, err := f.Stat()
	if err != nil {
		return false
	}

	return (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}
