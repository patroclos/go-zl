package search

import (
	"regexp"

	"jensch.works/zl/pkg/zettel"
)

type SearchQuery struct {
	Title *TextSearch
	Content *TextSearch
	Patterns SearchPatterns
}

// TODO: research how elastic or seq(sql like or dsl) model it
type TextSearch struct {
	Regex *regexp.Regexp
	Fuzzy *string
	Strict *string
}

type SearchPatterns []SearchPattern

// TODO: look at cypher, maybe implement in go
type SearchPattern struct {
	Topo TopoPattern
}

type TopoPattern struct {
	nodes []TopoSpec
}

type TopoSpec string

type Searcher interface {
	Search(query SearchQuery) (zettel.Id, error)
}
