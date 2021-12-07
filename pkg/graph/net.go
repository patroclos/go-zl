package graph

import "jensch.works/zl/pkg/zettel"

type Net struct {
	Nodes []zettel.Zettel
	Links []Link
}

type Node interface {
	Zettel() zettel.Zettel
	Refs() NodeRefs
}

type NodeRefs interface {
	All() []*Node
}

type Link struct {
	Source zettel.Zettel
	Target zettel.Zettel
}

type Spinner interface {
}
