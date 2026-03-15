package firebase

import (
	"context"
	"log"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var (
	app  *firebase.App
	once sync.Once
)

func Init() {
	once.Do(func() {
		credPath := os.Getenv("FIREBASE_CREDENTIALS")
		if credPath == "" {
			log.Fatal("FIREBASE_CREDENTIALS is required")
		}

		credJSON, err := os.ReadFile(credPath)
		if err != nil {
			log.Fatalf("failed to read firebase credentials file: %v", err)
		}

		app, err = firebase.NewApp(
			context.Background(),
			nil,
			option.WithCredentialsJSON(credJSON),
		)
		if err != nil {
			log.Fatalf("failed to initialize firebase: %v", err)
		}
	})
}

func App() *firebase.App {
	return app
}
