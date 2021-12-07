package main

import (
	"fmt"

	"jensch.works/zl/pkg/zettel"
)

func makeZ(zl zettel.Zettel) (z, error) {
	meta, err := zl.Metadata()
	if err != nil {
		return z{}, err
	}
	return z{
		Id:      string(zl.Id()),
		Title:   zl.Title(),
		Meta:    meta,
		TxtHref: fmt.Sprintf("https://localhost:8087/zettel/%s/text", zl.Id()),
	}, nil
}
