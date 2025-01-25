package repository

import (
	"maze-conquest-api/helper"
	"maze-conquest-api/model/domain"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MixStatsRepositoryImpl struct {
	FbApp *firebase.App
}

func NewMixStatsRepository(fbApp *firebase.App) MixStatsRepository {
	return &MixStatsRepositoryImpl{
		FbApp: fbApp,
	}
}

func (repository *MixStatsRepositoryImpl) GetMixStats(ctx *fiber.Ctx, uid string) *domain.MixStats {
	client, err := repository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc, err := client.Collection("mix_stats").Doc(uid).Get(ctx.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Handle the case where document doesn't exist
			emptyMixStats := new(domain.MixStats)
			emptyMixStats.Power = 0
			return emptyMixStats
		}

		// Throw panic for server error
		fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	data := doc.Data()
	userMixStats := helper.NewMixStats(data)
	userMixStats.Uid = doc.Ref.ID

	return userMixStats
}

func (repository *MixStatsRepositoryImpl) UpdatePower(ctx *fiber.Ctx, uid string, newPower int) bool {
	client, err := repository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Every time we update mix_stats, we update the user profile as well
	authClient, err := repository.FbApp.Auth(ctx.Context())
	if err != nil {
		panic(err)
	}
	user, err := authClient.GetUser(ctx.Context(), uid)
	helper.PanicForGetAuth(err, uid)

	doc := client.Collection("mix_stats").Doc(uid)
	_, err = doc.Update(ctx.Context(), []firestore.Update{
		{
			Path:  "power",
			Value: newPower,
		},
		{
			Path:  "ownerUsername",
			Value: user.DisplayName,
		},
		{
			Path:  "photoUrl",
			Value: user.PhotoURL,
		},
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Handle the case where document doesn't exist, Do not panic -> create new document
			// If everything is ok, no need to set uid since it's assign to key document
			mixStats := domain.MixStats{
				Power:         newPower,
				OwnerUsername: user.DisplayName,
				PhotoUrl:      user.PhotoURL,
			}
			client.Collection("mix_stats").Doc(uid).Set(ctx.Context(), mixStats)
			return true
		}

		// Throw panic for server error
		fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	return true
}

func (repository *MixStatsRepositoryImpl) GetLeaderboard(ctx *fiber.Ctx, uidCursor string) []*domain.MixStats {
	client, err := repository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	query := client.Collection("mix_stats").
		OrderBy("power", firestore.Desc).
		OrderBy(firestore.DocumentID, firestore.Desc). // Add consistent secondary ordering
		Limit(10)

	// Apply cursor if provided
	if uidCursor != "" {
		// Get the reference document first
		docSnap, err := client.Collection("mix_stats").Doc(uidCursor).Get(ctx.Context())
		if err != nil {
			panic(err)
		}

		query = query.StartAfter(docSnap)
	}

	// Execute the query
	docs, err := query.Documents(ctx.Context()).GetAll()
	if err != nil {
		panic(err)
	}

	var mixStatsSlice []*domain.MixStats
	for _, doc := range docs {
		mixStats := new(domain.MixStats)
		if err := doc.DataTo(mixStats); err != nil {
			panic(err)
		}
		mixStats.Uid = doc.Ref.ID
		mixStatsSlice = append(mixStatsSlice, mixStats)
	}

	return mixStatsSlice
}

func (repository *MixStatsRepositoryImpl) GetFirebaseInstance() *firebase.App {
	return repository.FbApp
}
