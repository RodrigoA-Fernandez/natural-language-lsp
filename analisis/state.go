package analisis

import (
	"log"
	"natural_language_lsp/lsp"
	"natural_language_lsp/scrapper"
	"strings"
	"unicode"
)

type State struct {
	// Map of Filenames to contents
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(document, text string) {
	s.Documents[document] = text
}

func (s *State) UpdateDocument(document, text string) {
	s.Documents[document] = text
}

func (s *State) Hover(id int, uri string, position lsp.Position, logger *log.Logger) lsp.HoverResponse {
	document := s.Documents[uri]
	documentSlice := strings.Split(document, "\n")
	linea := []rune(documentSlice[position.Line])

	start := position.Character
	for start > 0 {
		if !unicode.IsLetter(rune(linea[start])) {
			start++
			break
		}
		start--
	}

	end := position.Character
	for end < len(linea) {
		if !unicode.IsLetter(rune(linea[end])) {
			break
		}
		end++
	}

	palabra := linea[start:end]

	var texto string
	definicion, err := scrapper.Definir(string(palabra))
	if err != nil {
		texto = err.Error()
	} else {
		logger.Println(definicion)
		texto = scrapper.DefinicionMd(definicion)
	}

	response := lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			Id:  &id,
		},
		Result: lsp.HoverResult{
			Contents: lsp.MarkupContent{
				Kind:  lsp.Markdown,
				Value: texto,
			},
		},
	}
	return response
}
