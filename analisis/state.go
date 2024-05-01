package analisis

import (
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

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := s.Documents[uri]
	documentSlice := strings.Split(document, "\n")
	linea := documentSlice[position.Line]

	start := position.Character
	for start > 0 {
		if !unicode.IsLetter(rune(linea[start])){
			start++
			break
		}
		start--
	}

  end := position.Character
	for end < len(linea)-1 {
		if !unicode.IsLetter(rune(linea[end])){
			break
		}
		end++
	}

	palabra := linea[start:end]

	var texto string
	definicion, err := scrapper.Definir(palabra)
	if err != nil {
		texto = err.Error()
	} else {
		texto = scrapper.DefinicionMd(definicion)
	}

  texto += "\n" + palabra

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
