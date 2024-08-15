package grammarcheck

import (
	"log"
	"natural_language_lsp/analisis"
	"strings"

	katex "github.com/FurqanSoftware/goldmark-katex"
	lt "github.com/bas24/languagetool"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

func ParseDoc(uri string, logger *log.Logger, estado analisis.State) {
	md := goldmark.New(goldmark.WithExtensions(&katex.Extender{}, extension.TaskList))
	parser := md.Parser()
	node := parser.Parse(text.NewReader([]byte(estado.Documents[uri].Contenido)))

	// text := `Texto eroneo.`

	var textos []ast.Node
	getText(&textos, node, logger)

	arrayTextos := nodesToStringArray(textos, estado.Documents[uri].Contenido)
	textosCambiados := estado.GetChangedTexts(uri, arrayTextos, logger)

	for _, i := range textosCambiados {
		result, err := lt.Check(arrayTextos[i], "es-ES")
		if err != nil {
			logger.Println(err)
		}

		logger.Println(result)
	}

	// for _, v := range textos {
	// 	logger.Println(string(v.Text([]byte(doc))))
	// }
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

func nodesToStringArray(textos []ast.Node, doc string) []string {
	arrayTextos := []string{}
	for _, v := range textos {
		arrayTextos = append(arrayTextos, string(v.Text([]byte(doc))))
	}

	return arrayTextos
}
