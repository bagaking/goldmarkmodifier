package goldmarkmodifier

import (
	"github.com/yuin/goldmark/ast"
)

type (
	Mapper struct {
		Matcher
		Replacer
	}

	Matcher  func(source []byte, node ast.Node) bool
	Replacer func(source []byte, node ast.Node) []ast.Node
)

func NewMapper(matcher Matcher, replacer Replacer) Mapper {
	return Mapper{
		Matcher:  matcher,
		Replacer: replacer,
	}
}

func NewRemover(matcher Matcher) Mapper {
	return Mapper{
		Matcher: matcher,
		Replacer: func(source []byte, node ast.Node) []ast.Node {
			return nil
		},
	}
}
