package repository

import (
	"fmt"
	"maze-conquest-api/exception"
	"maze-conquest-api/helper"
	"maze-conquest-api/model/domain"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestRepositoryImpl struct {
	FbApp *firebase.App
}

func NewTestRepository(fbApp *firebase.App) TestRepository {
	return &TestRepositoryImpl{
		FbApp: fbApp,
	}
}

func (testRepository *TestRepositoryImpl) FindById(ctx *fiber.Ctx, uuid string) *domain.User {
	client, err := testRepository.FbApp.Firestore(ctx.Context())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	doc, err := client.Collection("users").Doc(uuid).Get(ctx.Context())
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Handle the case where document doesn't exist
			fmt.Println("user not found: " + err.Error())
			panic(exception.NewNotFoundError("User with ID '" + uuid + "' not found"))
		}

		// Throw panic for server error
		fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	data := doc.Data()
	data["uid"] = uuid

	user := helper.NewUser(data)

	return user
}

func (testRepository *TestRepositoryImpl) FindAll(ctx *fiber.Ctx) []*domain.User {
	client, err := testRepository.FbApp.Firestore(ctx.Context())
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
