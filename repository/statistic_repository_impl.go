package repository

import (
	"fmt"
	"maze-conquest-api/exception"
	"maze-conquest-api/model/domain"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatisticRepositoryImpl struct {
	FbApp *firebase.App
}

func NewStatisticRepository(fbApp *firebase.App) StatisticRepository {
	return &StatisticRepositoryImpl{
		FbApp: fbApp,
	}
}

func (repository *StatisticRepositoryImpl) GetUsers(ctx *fiber.Ctx) []*domain.Statistic {
	client, err := repository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc, err := client.Collection("statistics").Doc("users").Get(ctx.Context())
	if err != nil {
		// Throw panic for server error
		fiber.NewError(fiber.StatusInternalServerError, "Unexpected error")
	}

	data := doc.Data()
	worlds, ok := data["worlds"].(map[string]interface{})
	if !ok {
		panic("Error: 'worlds' is not a valid map[string]interface{}")
	}

	var statistics []*domain.Statistic
	for key, value := range worlds {
		// Check the type of value (e.g., assume it should be an integer for statistics)
		intValue, ok := value.(int64)
		if !ok {
			fmt.Printf("Error: key '%s' as it does not have an integer value\n", key)
			panic(fmt.Sprintf("Error: key '%s' as it does not have an integer value", key))
		}

		// Append each key-value pair to the statistics slice
		statistics = append(statistics, &domain.Statistic{
			Label: key,
			Value: float64(intValue),
		})
	}

	return statistics
}

func (repository *StatisticRepositoryImpl) GetMixStats(ctx *fiber.Ctx) []*domain.Statistic {
	client, err := repository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc, err := client.Collection("statistics").Doc("mix_stats").Get(ctx.Context())
	if err != nil {
		// Throw panic for server error
		fiber.NewError(fiber.StatusInternalServerError, "Unexpected error")
	}

	data := doc.Data()

	var statistics []*domain.Statistic
	for key, value := range data {
		// Check the type of value (e.g., assume it should be an integer for statistics)
		intValue, ok := value.(float64)
		if !ok {
			fmt.Printf("Error: key '%s' as it does not have an integer value\n", key)
			panic(fmt.Sprintf("Error: key '%s' as it does not have an integer value", key))
		}

		// Append each key-value pair to the statistics slice
		statistics = append(statistics, &domain.Statistic{
			Label: key,
			Value: float64(intValue),
		})
	}

	return statistics
}

func (repository *StatisticRepositoryImpl) GetUserPercentileFromLevel(ctx *fiber.Ctx, uid string) []*domain.Statistic {
	// 01 Requested user percentile, from 'power' document
	client, err := repository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc, err := client.Collection("statistics").Doc("users").Collection(uid).Doc("level").Get(ctx.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Handle the case where document doesn't exist

			// Handle when the User account are there but the data not yet created
			client, _ := repository.FbApp.Auth(ctx.Context())
			_, err := client.GetUser(ctx.Context(), uid)
			if err != nil {
				panic(exception.NewNotFoundError("User with ID '" + uid + "' not found"))
			}

			panic(exception.NewNotFoundError("Data not yet created"))
		}

		// Throw panic for server error
		panic("Failed to fetch user")
	}

	data := doc.Data()
	powerValue, ok := data["percentile"].(float64)
	if !ok {
		panic("Error: data['percentile'] as it does not have an float64 value")
	}

	// 02 Get global statistic for user model
	globalUserStatistic := repository.GetUsers(ctx)

	// 03 Combine
	merged := append(globalUserStatistic, &domain.Statistic{Value: powerValue, Label: "user"})
	return merged
}

func (repository *StatisticRepositoryImpl) GetUserPercentileFromPower(ctx *fiber.Ctx, uid string) []*domain.Statistic {
	// 01 Requested user percentile, from 'power' document
	client, err := repository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc, err := client.Collection("statistics").Doc("mix_stats").Collection(uid).Doc("power").Get(ctx.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Handle the case where document doesn't exist

			// Handle when the User account are there but the data not yet created
			client, _ := repository.FbApp.Auth(ctx.Context())
			_, err := client.GetUser(ctx.Context(), uid)
			if err != nil {
				panic(exception.NewNotFoundError("User with ID '" + uid + "' not found"))
			}

			panic(exception.NewNotFoundError("Data not yet created"))
		}

		// Throw panic for server error
		panic("Failed to fetch user")
	}

	data := doc.Data()
	powerValue, ok := data["percentile"].(float64)
	if !ok {
		panic("Error: data['percentile'] as it does not have an float64 value")
	}

	// 02 Get global statistic for user model
	globalUserStatistic := repository.GetMixStats(ctx)

	// 03 Combine
	merged := append(globalUserStatistic, &domain.Statistic{Value: powerValue, Label: "user"})
	return merged
}
