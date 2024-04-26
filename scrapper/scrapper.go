package scrapper

import (
	"errors"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Definicion struct {
	Palabra                  string
	Etimologia               string
	Acepciones               []string
	Definiciones_secundarias []Definicion
}

func Definir(palabra string) (Definicion, error) {
	var err error
	def := Definicion{}

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("user-agent", "Mozilla/5.0")
	})

	c.OnError(func(r *colly.Response, er error) {
		err = errors.New("No se ha podido conectar al diccionario.")
	})

	c.OnHTML("article", func(h *colly.HTMLElement) {
		def.Palabra = h.ChildText("header")

		def.Etimologia = h.ChildText(".n2")
		goquerySelection := h.DOM.Find(".j")
		goquerySelection.Each(func(i int, s *goquery.Selection) {
			cad := ""
			s.Each(func(i int, s *goquery.Selection) {
				cad = cad + " " + s.Text()
			})
			def.Acepciones = append(def.Acepciones, cad)
		})

		var frases_hechas []string
		var definiciones []string

		h.DOM.Find(".k5").Each(func(i int, s *goquery.Selection) {
			cad := ""
			s.Each(func(i int, s *goquery.Selection) {
				cad = cad + " " + s.Text()
			})
			frases_hechas = append(frases_hechas, cad)
		})

		h.DOM.Find(".k6").Each(func(i int, s *goquery.Selection) {
			cad := ""
			s.Each(func(i int, s *goquery.Selection) {
				cad = cad + " " + s.Text()
			})
			frases_hechas = append(frases_hechas, cad)
		})

		h.DOM.Find(".m").Each(func(i int, s *goquery.Selection) {
			cad := ""
			s.Each(func(i int, s *goquery.Selection) {
				cad = cad + " " + s.Text()
			})
			definiciones = append(definiciones, cad)
		})
		fmt.Println(definiciones)

		for i := 0; i < len(frases_hechas); i++ {
			var defs []string
			defs = append(defs, definiciones[i])
			def.Definiciones_secundarias = append(def.Definiciones_secundarias, Definicion{
				Palabra:    frases_hechas[i],
				Acepciones: defs,
			})
		}
	})

	c.Visit("https://dle.rae.es/js/20231220.js")
	c.Visit("https://dle.rae.es/" + palabra)
	if len(def.Palabra) == 0 && err == nil {
		err = errors.New("Palabra no encontrada en el Diccionario de la Lengua Española.")
	}

	return def, err
}
