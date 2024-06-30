package services

import (
	"context"
	"fmt"
	"log"
	"time"

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

func GetPriceAnalysisByItemId(ctx context.Context, itemId string) (*firestore.DocumentSnapshot, error) {
	firestoreClient, _ := clients.GetFirestoreClient()
	docs, firestoreError := firestoreClient.Collection(constants.PRICE_ANALYSIS_COLLECTION).Where("item_id", "==", itemId).Limit(1).Documents(ctx).GetAll()

	if firestoreError != nil {
		log.Fatalf("❌ Failed getting PRICE_ANALYSIS_COLLECTION element by ID | Error: %v", firestoreError)
	}

	if len(docs) > 0 {
		return docs[0], nil
	} else {
		return nil, fmt.Errorf("❌ No document found PRICE_ANALYSIS_COLLECTION with ID = %s", itemId)
	}
}

// =-=-==-=-=-=-=-=-=-=-=-=- METODOS UPDATE =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// La idea es usar este metodo para hacer la actualización masiva a la colección "price_analysis".
func UpdateItemBulk(ctx context.Context, itemsToUpdate []ScrapedItem) {
	// Obtenemos el cliente de firestore
	firestoreClient, _ := clients.GetFirestoreClient()

	// Creamos el cliente para hacer actualizacion bulk
	log.Printf("::: Creating BulkWrite client :::")
	bulkWrite := firestoreClient.BulkWriter(ctx)

	// Recorremos cada item a actualizar
	log.Printf("::: About to update %d items :::", len(itemsToUpdate))
	for i := 0; i < len(itemsToUpdate); i++ {
		item := itemsToUpdate[i]
		newMesurement := itemsToUpdate[i].ScrapResults

		// Buscamos en la BD el price_analysis relacionado
		log.Printf("::: About to GetPriceAnalysisByItemId [%s] :::", item.ItemId)
		priceAnalysisDoc, priceAnalysisError := GetPriceAnalysisByItemId(ctx, item.ItemId)

		// Si tiene price_analysis, actualizamos "measurements" y "last_analysis"
		if priceAnalysisError == nil {
			log.Printf("::: Previous price_analysis found :::")

			updatedMeasurements := mergePreviousMeasurementsWithNewOnes(priceAnalysisDoc, newMesurement)

			_, err := bulkWrite.Update(priceAnalysisDoc.Ref, []firestore.Update{
				{Path: "measurements", Value: updatedMeasurements},
				{Path: "last_analysis", Value: time.Now().UTC().Truncate(time.Second).String()},
			})

			if err != nil {
				// Handle any errors in an appropriate way, such as returning them.
				log.Printf("::: An error has occurred with item %s: %s :::", itemsToUpdate[i].ItemId, err)
			}
		} else {
			// Si NO tiene price_analysis, creamos un doc nuevo
			log.Printf("::: No previous price_analysis found, creating one :::")
			documentRef := firestoreClient.Collection(constants.PRICE_ANALYSIS_COLLECTION).NewDoc()

			_, err := bulkWrite.Create(documentRef, structures.Price_analysis{
				Measurements:              newMesurement,
				Last_analysis:             time.Now().UTC().Truncate(time.Second).String(),
				Item_id:                   item.ItemId,
				Analysis_interval_in_days: 7,
			})

			if err != nil {
				// y si explota lo decimos (?)
				log.Printf("::: An error has occurred with item %s: %s :::", itemsToUpdate[i].ItemId, err)
			}
		}

	}

	bulkWrite.Flush()

	log.Printf("::: Updated measurements! ✨ :::")
}

func mergePreviousMeasurementsWithNewOnes(priceAnalysisDoc *firestore.DocumentSnapshot, newMesurement map[string][]int) map[string][]int {
	var priceAnalysis structures.Price_analysis
	erroMapping := priceAnalysisDoc.DataTo(&priceAnalysis)
	if erroMapping != nil {
		log.Fatalf("❌ Failed mapping Price_analysis | Error: %v", erroMapping)
	}

	result := priceAnalysis.Measurements

	// Para cada "tienda", buscamos el historico y le agregamos el nuevo dato
	for key := range result {
		if newMesurement[key] != nil {
			result[key] = append(result[key], newMesurement[key][0])
		}
	}

	return result
}
