package grammarcheck

import (
	"log"
	"strings"

	"github.com/FurqanSoftware/goldmark-katex"
	// lt "github.com/bas24/languagetool"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

func ParseDoc(doc string, logger *log.Logger, textosOld []ast.Node) {
	md := goldmark.New(goldmark.WithExtensions(&katex.Extender{}, extension.TaskList))
	parser := md.Parser()
	node := parser.Parse(text.NewReader([]byte(doc)))
	// logger.Println(printTree(node, 0))
	// text := `Texto eroneo.`
	// result, err := lt.Check(text, "es-ES")
	// if err != nil {
	// 	logger.Println(err)
	// }

	// logger.Println(result)

	var textos []ast.Node
	getText(&textos, node, logger)
	for _, v := range textos {
		logger.Println(string(v.Text([]byte(doc))))
	}
}

func printTree(node ast.Node, indent int) string {
	cad := "\n" + strings.Repeat("\t", indent) + "(" + node.Kind().String()
	hijo := node.FirstChild()
	for hijo != nil {
		cad += printTree(hijo, indent+1)
		hijo = hijo.NextSibling()
	}
	cad += ")"
	return cad
}

func getText(textos *[]ast.Node, node ast.Node, logger *log.Logger) {
	if node.Kind() == ast.KindText {
		*textos = append(*textos, node)
	}
	hijo := node.FirstChild()
	for hijo != nil {
		getText(textos, hijo, logger)
		hijo = hijo.NextSibling()
	}
}
