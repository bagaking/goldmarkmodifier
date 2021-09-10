package goldmarkmodifier

import (
	"github.com/yuin/goldmark/ast"
)

type (
	Mapper struct {
		Matcher
		Replacer
	}

	Matcher  func(node ast.Node) bool
	Replacer func(node ast.Node) []ast.Node
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
		Replacer: func(node ast.Node) []ast.Node {
			return nil
		},
	}
}

var ParaRawTextCleaner = Mapper{
	Matcher: func(node ast.Node) bool {
		_, ok := node.(*ast.Paragraph)
		return ok
	},
	Replacer: func(node ast.Node) []ast.Node {
		n := ast.NewParagraph()
		allChild := make([]ast.Node, 0, n.ChildCount())
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			allChild = append(allChild, child)
		}
		for _, c := range allChild {
			n.AppendChild(n, c)
		}
		return []ast.Node{n}
	},
}
