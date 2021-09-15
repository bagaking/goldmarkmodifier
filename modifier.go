package goldmarkmodifier

import (
	"errors"
	"fmt"
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type (
	Modifier struct {
		node ast.Node

		ref *Modifier

		md     goldmark.Markdown
		source []byte
	}
)

func CreateNodeAndModifierBySource(md goldmark.Markdown, source []byte) (*Modifier, error) {
	node := md.Parser().Parse(text.NewReader(source))
	return CreateModifierBySourceAndNode(md, source, node)
}

func CreateModifierBySourceAndNode(md goldmark.Markdown, source []byte, node ast.Node) (*Modifier, error) {
	if node == nil {
		return nil, errors.New("node cannot by nil")
	}
	if source == nil {
		return nil, errors.New("source cannot be empty")
	}
	return &Modifier{
		source: source, node: node, md: md,
	}, nil
}

func (mod *Modifier) Root() *Modifier {
	if mod.ref == nil {
		return mod
	}
	return mod.ref.Root()
}

func (mod *Modifier) Node() ast.Node {
	return mod.node
}

func (mod *Modifier) Source() []byte {
	return mod.Root().source
}

func (mod *Modifier) Markdown() goldmark.Markdown {
	return mod.md
}

func (mod *Modifier) Render(w io.Writer) error {
	return mod.Markdown().Renderer().Render(w, mod.Source(), mod.Node())
}

// todo: using hash code to determine whether a node is belong to the modifier tree?
func (mod *Modifier) insertSource(source []byte) (from, to int) {
	root := mod.Root()
	orgLen := len(root.source)
	from, to = orgLen, orgLen+len(source)
	if from < to {
		root.source = append(root.source, source...)
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

func (mod *Modifier) WrapNode(source []byte) (ast.Node, error) {
	node := mod.Markdown().Parser().Parse(text.NewReader(source))
	from, _ := mod.insertSource(source)

	err := ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
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
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (mod *Modifier) WrapModifier(source []byte) (*Modifier, error) {
	node, err := mod.WrapNode(source)
	if err != nil {
		return nil, err
	}
	return mod.CreateSubNodeModifier(node)
}

func (mod *Modifier) CreateSubNodeModifier(node ast.Node) (*Modifier, error) {
	m := &Modifier{
		node: node, ref: mod.Root(),
	}
	m.ref = mod
	return m, nil
}

func (mod *Modifier) Dump() {
	fmt.Println("=== DUMP BEGIN ===")
	mod.Node().Dump(mod.source, 2)
	fmt.Println("=== DUMP END ===")
}

func (mod *Modifier) ReplaceNode(rs ...Mapper) error {
	return ast.Walk(mod.Node(), func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		for _, r := range rs {
			if !r.Matcher(mod.Source(), node) {
				continue
			}

			newNodes := r.Replacer(mod.Source(), node)
			parent := node.Parent()
			if parent == nil {
				continue
			}

			if len(newNodes) == 1 && newNodes[0] == node {
				break // already be modified
			}

			nextSibling := node.NextSibling()
			for _, n := range newNodes {
				if n == node {
					panic("you cannot just modify the node and insert it in a set") // which may cause the outter recursive in walk failed
				}
				parent.InsertBefore(parent, nextSibling, n)
			}
			parent.RemoveChild(parent, node)
			node.SetNextSibling(nextSibling)
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})
}
