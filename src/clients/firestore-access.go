package clients

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var (
	FirestoreClient *firestore.Client
	firestoreError  error
	once            sync.Once
)

func GetFirestoreClient() (*firestore.Client, error) {
	once.Do(func() {
		// We initialize Firestore client //
		ctx := context.Background()

		// ::: For local development ::: //
		// serviceAccount := option.WithCredentialsFile("C:/Users/alexi/Downloads/russell-5412-9b0867d4d571.json")
		// app, errNewApp := firebase.NewApp(ctx, nil, serviceAccount)
		// ::: End local development ::: //

		// ::: For deployed apps ::: //
		app, errNewApp := firebase.NewApp(ctx, nil)
		// ::: End local apps ::: //

		if errNewApp != nil {
			firestoreError = errNewApp
			log.Fatalln(firestoreError)
		}

		client, errConnet := app.Firestore(ctx)
		if errConnet != nil {
			firestoreError = errConnet
			log.Fatalln(firestoreError)
		}

		FirestoreClient = client

	})

	return FirestoreClient, firestoreError
}
