package module

import (
	"context"
	"log"
	"os"

	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const PhotoUrlWhenCrash string = "https://firebasestorage.googleapis.com/v0/b/maze-conquest-api.firebasestorage.app/o/placeholder_when_crash.webp?alt=media&token=3ca45c57-9581-461f-b15a-d17ec055d973"

func InitFirebase() (*firebase.App, error) {
	/*
		keys.json is only used for development.

		In production, keys.json is not included. Instead, the app will use build-keys.json.

		keys.json and build-keys.json contain the same information.
	*/
	var opt option.ClientOption
	if os.Getenv("MODE") == "prod" {
		opt = option.WithCredentialsFile("./build-keys.json")
	} else {
		opt = option.WithCredentialsFile("./keys.json")
	}
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app, err
}

// listenDocument listens to a single document.
func ListenDocument(ctx *fiber.Ctx, collection string, client *firestore.Client) error {
	// Set the correct response headers for SSE
	ctx.Set("Content-Type", "text/event-stream")
	ctx.Set("Cache-Control", "no-cache")
	ctx.Set("Connection", "keep-alive")

	// 01
	// projectID := "project-id"
	// Сontext with timeout stops listening to changes.
	timeoutCtx, cancel := context.WithTimeout(ctx.Context(), 30*time.Second)
	go func() {
		<-ctx.Context().Done() // Trigger cancel when HTTP client disconnects
		cancel()
	}()

	// 02
	it := client.Collection(collection).Doc("testId").Snapshots(timeoutCtx)
	for {
		snap, err := it.Next()
		// DeadlineExceeded will be returned when ctx is cancelled.
		if status.Code(err) == codes.DeadlineExceeded {
			return nil
		}

		// Handle other errors.
		if err != nil {
			return fmt.Errorf("Snapshots.Next: %w", err)
		}

		// Handle document no longer existing.
		if !snap.Exists() {
			fmt.Fprintf(ctx, "Document no longer exists\n")
			return nil
		}

		// Write snapshot data to the response writer.
		fmt.Fprintf(ctx, "Received document snapshot %v\n", snap.Data())
		fmt.Println(snap.Data())
	}
}
