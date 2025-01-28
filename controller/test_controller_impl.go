package controller

import (
	"maze-conquest-api/model/domain"
	"maze-conquest-api/model/web"
	"maze-conquest-api/repository"

	"github.com/gofiber/fiber/v2"
)

type TestControllerImpl struct {
	TestRepository repository.TestRepository
}

func NewTestController(testRepository repository.TestRepository) TestController {
	return &TestControllerImpl{
		TestRepository: testRepository,
	}
}

func (controller *TestControllerImpl) Gateway(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"code":   200,
		"status": "OK",
		"data": fiber.Map{
			"message": "Welcome to API Gateway",
		},
	})
}

func (controller *TestControllerImpl) FindById(ctx *fiber.Ctx) error {
	params := ctx.AllParams()
	uuid := params["uid"]

	var user *domain.User = controller.TestRepository.FindById(ctx, uuid)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	}

	return ctx.Status(200).JSON(webResponse)
}

func (controller *TestControllerImpl) FindAll(ctx *fiber.Ctx) error {
	var users []*domain.User = controller.TestRepository.FindAll(ctx)
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   users,
	}

	return ctx.Status(200).JSON(webResponse)
}
