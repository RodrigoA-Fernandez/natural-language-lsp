package analisis

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"natural_language_lsp/lsp"
	"natural_language_lsp/scrapper"
	"strings"
	"unicode"
)

type Document struct {
	Contenido string
	Textos    map[string]string
}

type State struct {
	// Map of Filenames to contents
	Documents map[string]Document
}

func NewState() State {
	return State{Documents: map[string]Document{}}
}

func (s *State) OpenDocument(document, text string) {
	s.Documents[document] = Document{text, map[string]string{}}
}

func (s *State) UpdateDocument(document, text string) {
	if entry, ok := s.Documents[document]; ok {
		entry.Contenido = text
	}
}

func (s *State) Hover(id int, uri string, position lsp.Position, logger *log.Logger) (lsp.HoverResponse, error) {
	document := s.Documents[uri]
	documentSlice := strings.Split(document.Contenido, "\n")
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

	logger.Println("----")
	logger.Printf("%d, %d\n", start, end)

	if start > end {
		return lsp.HoverResponse{}, errors.New("Fallo al parsear la palabra.")
	}

	palabra := linea[start:end]
	logger.Println(palabra)

	var texto string
	definicion, err := scrapper.Definir(string(palabra))
	if err != nil {
		texto = err.Error()
	} else {
		// logger.Println(definicion)
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
	return response, nil
}

func (s *State) GetChangedTexts(texts []string) map[string]string {
	nuevoMapa := map[string]string{}
	textosCambiados := map[string]string{}

	for _, v := range texts {
		hash := hashMD5(v)
		_, ok := s.Documents[hash]
		if !ok {
			textosCambiados[hash] = v
		}
		nuevoMapa[hash] = v
	}
	return textosCambiados
}

func hashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
