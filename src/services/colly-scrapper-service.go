package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/exp/maps"
	"russell.com/hardware_scrapper/clients"
	"russell.com/hardware_scrapper/structures"
)

func ScrapItem(target structures.Items_in_target) map[string]int {
	collyClient := clients.GetCollyClient()
	targets := target.Targets
	storeNames := maps.Keys(targets)

	resultsMap := make(map[string]int)
	for _, store := range storeNames {
		resultsMap[store] = 0
	}

	for i := 0; i < len(targets); i++ {
		targetInfo := targets[storeNames[i]]

		collyClient.OnHTML(targetInfo.HtmlTarget, func(e *colly.HTMLElement) {
			priceFound := e.Text

			stringWithoutSymbol := strings.TrimSpace(strings.ReplaceAll(priceFound, "$", ""))

			intValue, err := strconv.Atoi(stringWithoutSymbol)
			if err != nil {
				fmt.Println("🦜 Error converting string to int:", err)
			} else {
				resultsMap[storeNames[i]] = intValue
			}

			// Matamos todo lo relacionado a la pagina que visitamos
			collyClient.OnHTMLDetach(targetInfo.HtmlTarget)

		})

		// Visitamos la url
		collyClient.Visit(targetInfo.Url)
	}

	return resultsMap

}
