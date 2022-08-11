package prompt

import (
	"git.jensch.dev/joshua/go-zl/pkg/zettel/elemz"
)

// The prompt package is a parser for a small number of training/reminder question-types
// such as a Question-Answer format or Omitting holes from a Sentence/Paragraph.

func Parser() elemz.Parser {
	return elemz.NewCompoundParser(&parseQA{}, &parseOmit{})
}
