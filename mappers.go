package goldmarkmodifier

import (
	"github.com/yuin/goldmark/ast"
)

var mapperCleanRawText = Mapper{
	Matcher: func(source []byte, node ast.Node) bool {
		_, ok := node.(*ast.Paragraph)
		return ok
	},
	Replacer: func(source []byte, node ast.Node) []ast.Node {
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

func MCleanRawText() Mapper {
	return mapperCleanRawText
}

func MMoveAllHeaderLevel(matcher Matcher, move int, max int) Mapper {
	return NewMapper(func(source []byte, node ast.Node) bool {
		_, ok := node.(*ast.Heading)
		if !ok {
			return false
		}
		if matcher == nil {
			return true
		}
		return matcher(source, node)
	}, func(source []byte, node ast.Node) []ast.Node {
		tnHead := node.(*ast.Heading)
		tnHead.Level += move
		if tnHead.Level < 1 {
			tnHead.Level = 1
		} else if max != 0 && max > 1 && tnHead.Level > max {
			tnHead.Level = max
		}
		return []ast.Node{tnHead}
	})
}
