package graph

import (
	"jensch.works/zl/pkg/zettel"
)

// https://pkg.go.dev/gonum.org/v1/gonum/graph?utm_source=godoc#Graph

type Node struct {
	Z zettel.Z
}

func (n Node) ID() int64 {
	return int64(n.Z.Metadata().CreateTime.Nanosecond())
}
