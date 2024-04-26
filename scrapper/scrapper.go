package scrapper

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Definicion struct {
	Palabra               string
	Etimologia            string
	Acepciones            []string
	Definicion_secundaria []Definicion
}

func Definir(palabra string) (Definicion, error) {
	var err error
	def := Definicion{}

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("user-agent", "Mozilla/5.0")
	})

	c.OnError(func(r *colly.Response, er error) {
		// fmt.Println("Code:", r.StatusCode)
		// fmt.Println("Err:", err)
		// fmt.Println(r.Ctx)
		// fmt.Println(string(r.Body))
		// fmt.Println(r.Headers)
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
	})

	c.Visit("https://dle.rae.es/js/20231220.js")
	c.Visit("https://dle.rae.es/" + palabra)
	return def, err
}
