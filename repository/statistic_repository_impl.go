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

func (repository *StatisticRepositoryImpl) GetUserLeaderboard(ctx *fiber.Ctx, uid string) *domain.Leaderboard {
	client, err := repository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Create reference so we don't write the same code
	leaderboardDocRef := client.Collection("statistics").Doc("leaderboard")

	// Get the specific user info
	userDoc, err := leaderboardDocRef.Collection("users").Doc(uid).Get(ctx.Context())
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
		panic(fmt.Sprintf("Failed to fetch user; Unhandled error while get user leaderboard %s", uid))
	}

	// Get the global info
	leaderboardDoc, err := leaderboardDocRef.Get(ctx.Context())
	if err != nil {
		panic(err)
	}

	// Get the field data
	userData := userDoc.Data()
	leaderboardData := leaderboardDoc.Data()

	// Extract the data
	/*
		fmt.Println(reflect.TypeOf(userData["percentile"]), "userData['percentile']")
		fmt.Println(reflect.TypeOf(userData["rank"]), "userData['rank']")
		fmt.Println(reflect.TypeOf(userData["total"]), "userData['total']")
		fmt.Println(reflect.TypeOf(leaderboardData["average"]), "leaderboardData['average']")
		fmt.Println(reflect.TypeOf(leaderboardData["len"]), "leaderboardData['len']")
	*/

	percentile, ok := userData["percentile"].(float64)
	if !ok {
		panic("Error while check the .percentile; type not float64")
	}

	rankFloat, ok := userData["rank"].(int64)
	if !ok {
		panic("Error while check .rank; type not int64")
	}
	rank := int(rankFloat)

	total, ok := userData["total"].(float64)
	if !ok {
		panic("Error while check .total; type not float64")
	}

	var average float64
	switch v := leaderboardData["average"].(type) {
	case int:
		average = float64(v) // Convert int to float64
	case int64:
		average = float64(v) // Convert int64 to float64
	case float64:
		average = v // Already float64, no conversion needed
	default:
		panic(fmt.Sprintf("unexpected type for average: %T", v))
	}

	leaderboardLength64, ok := leaderboardData["len"].(int64)
	if !ok {
		panic("Error while check .len; ;type not int64")
	}
	leaderboardLength := int(leaderboardLength64)

	leaderboard := &domain.Leaderboard{
		GlobalAverage:  average,
		TotalUser:      leaderboardLength,
		UserPercentile: percentile,
		UserRank:       rank,
		UserTotalPower: total,
	}

	return leaderboard
}
