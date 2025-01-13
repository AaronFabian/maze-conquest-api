package repository

import (
	"fmt"
	"maze-conquest-api/exception"
	"maze-conquest-api/helper"
	"maze-conquest-api/model/domain"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRepositoryImpl struct {
	FbApp *firebase.App
}

func NewUserRepository(fbApp *firebase.App) UserRepository {
	return &UserRepositoryImpl{
		FbApp: fbApp,
	}
}

func (userRepository *UserRepositoryImpl) FindById(ctx *fiber.Ctx, uid string) *domain.User {
	client, err := userRepository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc, err := client.Collection("users").Doc(uid).Get(ctx.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Handle the case where document doesn't exist
			fmt.Println("user not found: " + err.Error())
			panic(exception.NewNotFoundError("User with ID '" + uid + "' not found"))
		}

		// Throw panic for server error
		fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	data := doc.Data()
	data["uid"] = uid

	user := helper.NewUser(data)

	return user
}

func (userRepository *UserRepositoryImpl) FindAll(ctx *fiber.Ctx) []*domain.User {
	client, err := userRepository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	docs, err := client.Collection("users").Documents(ctx.Context()).GetAll()
	if err != nil {
		panic(err)
	}

	var users []*domain.User
	for _, doc := range docs {
		var user domain.User
		if err := doc.DataTo(&user); err != nil {
			panic(err)
		}
		user.Uid = doc.Ref.ID
		users = append(users, &user)
	}

	return users
}

func (userRepository *UserRepositoryImpl) UpdateItem(ctx *fiber.Ctx, uid string, itemName string, quantity int) *domain.User {
	client, err := userRepository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc := client.Doc("users/" + uid)
	_, err = doc.Update(ctx.Context(), []firestore.Update{
		{
			Path:  "items." + itemName,
			Value: quantity,
		},
	})
	if err != nil {
		panic(err)
	}

	//
	newDoc, err := doc.Get(ctx.Context())
	if err != nil {
		panic(err)
	}

	data := newDoc.Data()
	data["uid"] = uid

	user := helper.NewUser(data)

	return user
}

func (userRepository *UserRepositoryImpl) GetAllHeroes(ctx *fiber.Ctx, uid string) []*domain.Hero {
	client, err := userRepository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// test := client.Collection("users").Select("address").Doc(uid).Get(ctx)
	// iter, err := client.Collection("users").Doc("testId").Collection("allHeroes").Doc("soldier").Get(ctx.Context())
	// fmt.Println(iter.Data())

	doc, err := client.Collection("users").Doc(uid).Get(ctx.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Handle the case where document doesn't exist
			fmt.Println("user not found: " + err.Error())
			panic(exception.NewNotFoundError("User with ID '" + uid + "' not found"))
		}

		// Throw panic for server error
		fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	// Extract the allHeroes field (assuming it’s a map)
	allHeroes, ok := doc.Data()["allHeroes"].(map[string]interface{})
	if !ok {
		fmt.Println("[System] Fatal error while accessing .allHeroes")
		panic("Fatal error while accessing .allHeroes")
	}

	var heroes []*domain.Hero
	for hero, data := range allHeroes {
		heroData, ok := data.(map[string]interface{})
		if !ok {
			fmt.Println("[System] Fatal error while accessing .allHeroes.data")
			panic("Fatal error while accessing .allHeroes.data")
		}
		// for key, value := range heroData {
		// 	fmt.Printf("  %s: %v\n", key, value)
		// }

		level, ok := heroData["level"].(int64)
		if !ok {
			fmt.Println("[System] Fatal error, 'level' properties is not int64")
			panic("Fatal error, 'level' properties is not int64")
		}

		hero := &domain.Hero{
			Name:  hero,
			Level: int(level),
		}
		heroes = append(heroes, hero)
	}

	return heroes
}

func (userRepository *UserRepositoryImpl) MazeLevel(ctx *fiber.Ctx, uid string) *domain.World {
	client, err := userRepository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc, err := client.Collection("users").Doc(uid).Get(ctx.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			panic(exception.NewNotFoundError("User with ID '" + uid + "' not found"))
		}

		fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	worlds, ok := doc.Data()["worlds"].(map[string]interface{})
	if !ok {
		panic("Error while accessing worlds properties from user with Id " + uid)
	}

	// Only want the maze level
	worldName := "level"
	mazeLevel, _ := worlds[worldName].(int64)
	world := domain.World{
		Name:  worldName,
		Level: int(mazeLevel),
	}

	return &world
}
