package scrapper

import (
	"errors"
	"html"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Definicion struct {
	Palabra                  []string
	Etimologia               []string
	Acepciones               [][]string
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
		def.Palabra = append(def.Palabra, h.ChildText("header"))
		def.Etimologia = append(def.Etimologia, h.ChildText(".n2"))
		goquerySelection := h.DOM.Find(".j, .l2")

		var ac []string
		goquerySelection.Each(func(i int, s *goquery.Selection) {
			cad := ""
			s.Each(func(i int, s *goquery.Selection) {
				cad = cad + " " + s.Text()
			})
			ac = append(ac, cad)
		})
		def.Acepciones = append(def.Acepciones, ac)

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

		for i := 0; i < len(frases_hechas); i++ {
			var defs []string
			defs = append(defs, definiciones[i])
			def.Definiciones_secundarias = append(def.Definiciones_secundarias, Definicion{
				Palabra:    []string{frases_hechas[i]},
				Acepciones: [][]string{defs},
			})
		}
	})

	palabraHtml := html.EscapeString(palabra)

	c.Visit("https://dle.rae.es/js/20231220.js")
	c.Visit("https://dle.rae.es/" + palabraHtml)
	if len(def.Palabra) == 0 && err == nil {
		err = errors.New("Palabra no encontrada en el Diccionario de la Lengua Española.")
	}

	return def, err
}

func DefinicionMd(def Definicion) string {
	texto := ""
	for i, palabra := range def.Palabra {
		texto = texto + "# " + palabra + "\n"
		if len(def.Etimologia) != 0 {
			texto += "*" + def.Etimologia[i] + "*\n"
		}
		for _, acepcion := range def.Acepciones[i] {
			texto += acepcion + "\n"
		}
		texto += "\n"
	}
	for _, sec := range def.Definiciones_secundarias {
		texto += "## " + sec.Palabra[0] + "\n"
		texto += sec.Acepciones[0][0] + "\n\n"
	}
	return texto
}
