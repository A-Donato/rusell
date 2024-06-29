package clients

import (
	"fmt"
	"log"
	"sync"

	"github.com/gocolly/colly"
)

var (
	CollyClient *colly.Collector
	onceColly   sync.Once
)

func GetCollyClient() *colly.Collector {
	onceColly.Do(func() {

		c := colly.NewCollector()

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("🦜 solicitamos la pagina: ", r.URL)
		})

		c.OnError(func(_ *colly.Response, err error) {
			log.Println("🦜 Explotamos: ", err)
		})

		c.OnResponse(func(r *colly.Response) {
			fmt.Println("🦜 obtuvimos la pagina: ", r.Request.URL)
		})

		c.OnScraped(func(r *colly.Response) {
			fmt.Println("🦜 scraped!", r.Request.URL)
		})

		CollyClient = c
	})

	return CollyClient
}
