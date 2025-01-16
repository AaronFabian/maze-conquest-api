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

func NewMixStatsRepositoryImpl(fbApp *firebase.App) MixStatsRepository {
	return &MixStatsRepositoryImpl{
		FbApp: fbApp,
	}
}

func (mixStatsRepository *MixStatsRepositoryImpl) GetMixStats(ctx *fiber.Ctx, uid string) *domain.MixStats {
	client, err := mixStatsRepository.FbApp.Firestore(ctx.Context())
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

func (mixStatsRepository *MixStatsRepositoryImpl) UpdatePower(ctx *fiber.Ctx, uid string, newPower int) bool {
	client, err := mixStatsRepository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc := client.Collection("mix_stats").Doc(uid)
	_, err = doc.Update(ctx.Context(), []firestore.Update{
		{
			Path:  "power",
			Value: newPower,
		},
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Handle the case where document doesn't exist, Do not panic -> create new document
			mixStats := domain.MixStats{
				Power: newPower,
			}
			client.Collection("mix_stats").Doc(uid).Set(ctx.Context(), mixStats)
			return true
		}

		// Throw panic for server error
		fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	return true
}

func (mixStatsRepository *MixStatsRepositoryImpl) GetLeaderboard(ctx *fiber.Ctx, uidCursor string) []*domain.MixStats {
	client, err := mixStatsRepository.FbApp.Firestore(ctx.Context())
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
		var mixStats domain.MixStats
		if err := doc.DataTo(&mixStats); err != nil {
			panic(err)
		}
		mixStats.Uid = doc.Ref.ID
		mixStatsSlice = append(mixStatsSlice, &mixStats)
	}

	return mixStatsSlice
}

func (mixStatsRepository *MixStatsRepositoryImpl) GetFirebaseInstance() *firebase.App {
	return mixStatsRepository.FbApp
}
