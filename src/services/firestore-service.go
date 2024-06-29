package services

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"russell.com/hardware_scrapper/clients"
	"russell.com/hardware_scrapper/constants"
	"russell.com/hardware_scrapper/structures"
)

// =-=-==-=-=-=-=-=-=-=-=-=- METODOS GET =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

func GetItems(ctx context.Context) []structures.Item {
	firestoreClient, _ := clients.GetFirestoreClient()
	var result []structures.Item

	docs, error := firestoreClient.Collection(constants.ITEMS_COLLECTION).Limit(10).Documents(ctx).GetAll()

	if error != nil {
		log.Fatalf("❌ Failed getting ITEMS_COLLECTION | Error: %v", error)
	}

	log.Printf("estamos intentando transformar los items")

	for i := 0; i < len(docs); i++ {
		var item structures.Item

		erroMapping := docs[i].DataTo(&item)
		if erroMapping != nil {
			log.Fatalf("❌ Failed mapping ITEMS_COLLECTION | Error: %v", erroMapping)
		}

		result = append(result, item)
	}

	return result
}

func GetItemsInTarget(ctx context.Context) []structures.Items_in_target {
	firestoreClient, _ := clients.GetFirestoreClient()
	var result []structures.Items_in_target

	docs, error := firestoreClient.Collection(constants.ITEMS_IN_TARGET_COLLECTION).Limit(10).Documents(ctx).GetAll()

	if error != nil {
		log.Fatalf("❌ Failed getting ITEMS_IN_TARGET_COLLECTION | Error: %v", error)
	}

	log.Printf("estamos intentando transformar los items")

	for i := 0; i < len(docs); i++ {
		var item structures.Items_in_target

		erroMapping := docs[i].DataTo(&item)
		if erroMapping != nil {
			log.Fatalf("❌ Failed mapping ITEMS_IN_TARGET_COLLECTION | Error: %v", erroMapping)
		}

		result = append(result, item)
	}

	return result
}

func GetScrapTargets(ctx context.Context) []*firestore.DocumentSnapshot {
	firestoreClient, _ := clients.GetFirestoreClient()
	docs, error := firestoreClient.Collection(constants.SCRAP_TARGETS_COLLECTION).Limit(10).Documents(ctx).GetAll()

	if error != nil {
		log.Fatalf("❌ Failed getting SCRAP_TARGETS_COLLECTION | Error: %v", error)
	}

	return docs
}

func GetPriceAnalysis(ctx context.Context) []*firestore.DocumentSnapshot {
	firestoreClient, _ := clients.GetFirestoreClient()
	docs, error := firestoreClient.Collection(constants.PRICE_ANALYSIS_COLLECTION).Limit(10).Documents(ctx).GetAll()

	if error != nil {
		log.Fatalf("❌ Failed getting PRICE_ANALYSIS_COLLECTION | Error: %v", error)
	}

	return docs
}

// =-=-==-=-=-=-=-=-=-=-=-=- METODOS UPDATE =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-
type Item_update_payload struct {
	documentRef *firestore.DocumentRef
	Measurement int
	itemId      string
}

// La idea es usar este metodo para hacer la actualización masiva a la colección "price_analysis".
func UpdateItemBulk(ctx context.Context, itemsToUpdate []Item_update_payload) {
	// Obtenemos el cliente de firestore
	firestoreClient, _ := clients.GetFirestoreClient()

	// Get a new bulkWrite.
	log.Printf("::: Creating BulkWrite client :::")
	bulkWrite := firestoreClient.BulkWriter(ctx)

	log.Printf("::: About to update %d items :::", len(itemsToUpdate))
	for i := 0; i < len(itemsToUpdate); i++ {
		docRef := itemsToUpdate[i].documentRef
		newMesurement := itemsToUpdate[i].Measurement

		_, err := bulkWrite.Update(docRef, []firestore.Update{
			{Path: "measurements", Value: firestore.ArrayUnion(newMesurement)},
		})

		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("::: An error has occurred with item %s: %s :::", itemsToUpdate[i].itemId, err)
		}
	}

	bulkWrite.Flush()

	log.Printf("::: Updated measurements! ✨ :::")
}
