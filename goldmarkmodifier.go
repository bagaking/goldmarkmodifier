package goldmarkmodifier

import (
	"fmt"
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type (
	Modifier struct {
		md     goldmark.Markdown
		source []byte
		node   ast.Node
	}
)

func NewModifier(md goldmark.Markdown, source []byte) Modifier {
	node := md.Parser().Parse(text.NewReader(source))
	return Modifier{
		source: source, node: node, md: md,
	}
}

func (mod *Modifier) Root() ast.Node {
	return mod.node
}

func (mod *Modifier) Source() []byte {
	return mod.source
}

func (mod *Modifier) Render(w io.Writer) error {
	return mod.md.Renderer().Render(w, mod.source, mod.node)
}

func (mod *Modifier) insertSource(source []byte) (from, to int) {
	orgLen := len(mod.source)
	from, to = orgLen, orgLen+len(source)
	if from < to {
		mod.source = append(mod.source, source...)
	} else {
		to = from
	}
	return
}

func (mod *Modifier) WarpText(source string) ast.Node {
	nd := ast.NewText()
	from, to := mod.insertSource([]byte(source))

	nd.Segment = text.NewSegment(from, to)
	return nd
}

func (mod *Modifier) WrapNode(source []byte) ast.Node {
	node := mod.md.Parser().Parse(text.NewReader(source))
	from, _ := mod.insertSource(source)

	ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			return ast.WalkContinue, nil
		}
		switch tn := node.(type) {
		case *ast.Text:
			tn.Segment = text.NewSegment(tn.Segment.Start+from, tn.Segment.Stop+from)
		case *ast.RawHTML:
			for i := 0; i < tn.Segments.Len(); i++ {
				org := tn.Segments.At(i)
				tn.Segments.Set(i, text.NewSegment(org.Start+from, org.Stop+from))
			}
		}
		if node.Type() != ast.TypeInline {
			for i, lines := 0, node.Lines(); i < lines.Len(); i++ {
				org := lines.At(i)
				lines.Set(i, text.NewSegment(org.Start+from, org.Stop+from))
			}
		}
		return ast.WalkContinue, nil
	})
	return node
}

func (mod *Modifier) Dump() {
	fmt.Println("=== DUMP BEGIN ===")
	mod.Root().Dump(mod.source, 2)
	fmt.Println("=== DUMP END ===")
}

func (mod *Modifier) ReplaceNode(rs ...Mapper) {
	ast.Walk(mod.Root(), func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		for _, r := range rs {
			if !r.Matcher(node) {
				continue
			}

			newNodes := r.Replacer(node)
			parent, next := node.Parent(), node.NextSibling()
			if parent == nil {
				continue
			}

			cur := node
			for _, n := range newNodes {
				parent.InsertAfter(parent, cur, n)
				cur = n
			}

			parent.RemoveChild(parent, node)
			node.SetNextSibling(next)
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})
}
