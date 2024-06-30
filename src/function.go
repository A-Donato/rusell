package hardware_scrapper

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"russell.com/hardware_scrapper/clients"
	"russell.com/hardware_scrapper/services"
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
	// defer FirestoreClient.Close()

	// Obtenemos los items base
	// items := services.GetItems(ctx)

	// Obtenemos los items para analizar
	itemsInTarget := services.GetItemsInTarget(ctx)

	// Scrappeamos
	result := services.ScrapItem(itemsInTarget[0])
	log.Println("resultado de la ejecución: ", result)

	// Actualizamos el historico de precios
	services.UpdateItemBulk(ctx, []services.ScrapedItem{result})

	fmt.Fprintf(w, "no fuimos")

}
