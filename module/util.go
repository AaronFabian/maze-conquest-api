package module

import (
	"context"
	"log"
	"os"

	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const PhotoUrlWhenCrash string = "https://firebasestorage.googleapis.com/v0/b/maze-conquest-api.firebasestorage.app/o/placeholder_when_crash.webp?alt=media&token=3ca45c57-9581-461f-b15a-d17ec055d973"

func getSecretManager() []byte {
	projectID := os.Getenv("GCP_PROJECT_ID")
	secretID := os.Getenv("GCP_SECRET_ID")
	const version = "latest"

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	// Build the request name
	secretName := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, secretID, version)

	// Access the secret
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	data := result.Payload.Data
	return data
}

func InitFirebase() (*firebase.App, error) {
	sm := getSecretManager()
	opt := option.WithCredentialsJSON(sm)
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

func EnvConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, assuming using production env")
	}
}
