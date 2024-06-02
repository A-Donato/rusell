package hardware_scrapper

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gocolly/colly"
	"russell.com/hardware_scrapper/clients"
)

// defining a data structure to store the scraped data
type PokemonProduct struct {
	Url, Image, Name, Price string
}

var (
	FirestoreClient *firestore.Client
	ctx             context.Context
)

func init() {
	// Inicializamos variables globales
	FirestoreClient, _ = clients.GetFirestoreClient()
	ctx = context.Background()

	// Definimos todas las funciones de entrada
	functions.HTTP("start-scrapping", scrappHardware)
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
func scrappHardware(w http.ResponseWriter, r *http.Request) {
	// Cerramos el cliente de firebase
	defer FirestoreClient.Close()

	// Conectamos la db
	// Use a service account

	// Scrapper
	// initializing the slice of structs that will contain the scraped data
	var pokemonProducts []PokemonProduct

	// // scraping logic...
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	// iterating over the list of HTML product elements
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		// initializing a new PokemonProduct instance
		pokemonProduct := PokemonProduct{}

		// scraping the data of interest
		pokemonProduct.Url = e.ChildAttr("a", "href")
		pokemonProduct.Image = e.ChildAttr("img", "src")
		pokemonProduct.Name = e.ChildText("h2")
		pokemonProduct.Price = e.ChildText(".price")

		// adding the product instance with scraped data to the list of products
		pokemonProducts = append(pokemonProducts, pokemonProduct)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	// downloading the target HTML page
	c.Visit("https://scrapeme.live/shop/")

	// adding each Pokemon product to the CSV output file
	for _, pokemonProduct := range pokemonProducts {
		pokemon := map[string]interface{}{
			"url":   pokemonProduct.Url,
			"image": pokemonProduct.Image,
			"name":  pokemonProduct.Name,
			"price": pokemonProduct.Price,
		}

		_, err := FirestoreClient.Collection("pokemons").Doc(pokemonProduct.Name).Set(ctx, pokemon)

		if err != nil {
			log.Fatalf("Failed adding pokemon %v | Error: %v", pokemonProduct.Name, err)
		} else {
			log.Printf("Guardamos a: %v", pokemonProduct.Name)
		}
	}
	// log.Println("Terminamos")

	fmt.Fprintf(w, "Terminamos con la subida de mierda")
}
