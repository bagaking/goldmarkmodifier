package main

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/bagaking/goldmarkmodifier"
	"github.com/bagaking/gotools/strs"
)

var testFile = `
# Title

contents

[### xx](https://localhost)123
`

var anotherFile = `
##### another file

xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  
yyyyyyyyyyyyyyyyyyyyyyyyyyyyyy

`

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Table),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	mod := goldmarkmodifier.NewModifier(md, []byte(testFile), nil)
	mod.Dump()

	bDump := make([]byte, 0, 1000)
	fmt.Println("=== DUMP ===", string(bDump))

	para := ast.NewParagraph()
	head := ast.NewHeading(2)
	para.AppendChild(para, head)
	head.AppendChild(head, mod.WarpText("wrapped head"))
	mod.Root().AppendChild(mod.Root(), para)

	mod.ReplaceNode(goldmarkmodifier.MCleanRawText())
	mod.Dump()

	mapperInsert := goldmarkmodifier.NewMapper(func(source []byte, node ast.Node) bool {
		switch tn := node.(type) {
		case *ast.Link:
			text := tn.Text(source)
			if strs.StartsWith(string(text), "#") {
				return true
			}
		}
		return false
	}, func(source []byte, node ast.Node) []ast.Node {
		newNode := mod.WrapNode([]byte(anotherFile))
		aaaPara := ast.NewParagraph()
		aaaPara.AppendChild(aaaPara, newNode)
		mod.Root().AppendChild(mod.Root(), aaaPara)
		return []ast.Node{newNode}
	})

	mod.ReplaceNode(mapperInsert)
	mod.Dump()

	var buf bytes.Buffer
	if err := mod.Render(&buf); err != nil {
		panic(err)
	}
	fmt.Println("\n\nOUTPUT:\n\n", buf.String())
}
